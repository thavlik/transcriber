package server

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"github.com/thavlik/transcriber/pharmaseer/pkg/thumbcache"
	"go.uber.org/zap"
)

type entityType int

func (e entityType) String() string {
	switch e {
	case drugEntity:
		return "drug"
	default:
		panic(base.Unreachable)
	}
}

const (
	drugEntity entityType = 1
)

type entity struct {
	Type  entityType `json:"t"`
	Query string     `json:"q"`
	Force bool       `json:"f"`
}

func initQueryWorkers(
	concurrency int,
	infoCache infocache.InfoCache,
	pub pubsub.Publisher,
	querySched scheduler.Scheduler,
	pdbSched scheduler.Scheduler,
	svgCache thumbcache.ThumbCache,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popQuery := make(chan string)
	go queryPopper(popQuery, querySched, stop, log)
	for i := 0; i < concurrency; i++ {
		go queryWorker(
			infoCache,
			popQuery,
			querySched,
			pdbSched,
			svgCache,
			pub,
			log,
		)
	}
}

func queryPopper(
	popQuery chan<- string,
	querySched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := querySched.Notify()
	defer close(popQuery)
	delay := 12 * time.Second
	for {
		start := time.Now()
		entities, err := querySched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(entities) > 0 {
			log.Debug("checking entities", zap.Int("num", len(entities)))
			for _, entJson := range entities {
				popQuery <- entJson
			}
		}
		remaining := delay - time.Since(start)
		if remaining < 0 {
			remaining = time.Millisecond
		}
		select {
		case <-stop:
			return
		case <-notification:
			continue
		case <-time.After(remaining):
			continue
		}
	}
}

func queryWorker(
	infoCache infocache.InfoCache,
	popQuery <-chan string,
	querySched scheduler.Scheduler,
	pdbSched scheduler.Scheduler,
	svgCache thumbcache.ThumbCache,
	pub pubsub.Publisher,
	log *zap.Logger,
) {
	subscriptionKey, ok := os.LookupEnv("BING_API_KEY")
	if !ok {
		panic("Ocp-Apim-Subscription-Key not set")
	}
	for {
		rawEnt, ok := <-popQuery
		if !ok {
			return
		}
		e := &entity{}
		if err := json.Unmarshal([]byte(rawEnt), e); err != nil {
			log.Warn("removing malformed entity", zap.Error(err))
			if err := querySched.Remove(rawEnt); err != nil {
				log.Error("failed to remove malformed entity",
					zap.Error(err),
					zap.String("entity", rawEnt))
				panic(errors.Wrap(err, "failed to remove malformed entity"))
			}
			continue
		}
		entityLog := log.With(
			zap.String("type", e.Type.String()),
			zap.String("query", e.Query))
		lock, err := querySched.Lock(e.Query)
		if err == scheduler.ErrLocked {
			// go to the next project
			entityLog.Debug("entity already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		entityLog.Debug("processing drug query")
		if err := func() error {
			defer lock.Release()
			ctx, cancel := context.WithCancel(context.Background())
			onProgress := make(chan struct{}, 4)
			stopped := make(chan struct{})
			defer func() {
				cancel()
				<-stopped
			}()
			go func() {
				defer func() { stopped <- struct{}{} }()
				for {
					select {
					case <-ctx.Done():
						return
					case _, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							panic(err)
						}
					}
				}
			}()
			switch e.Type {
			case drugEntity:
				var drug *api.DrugDetails
				if !e.Force {
					drug, err = infoCache.GetDrug(
						ctx,
						e.Query,
					)
					if err != nil && err != infocache.ErrCacheUnavailable {
						return errors.Wrap(err, "failed to get drug from cache")
					}
				}
				if drug == nil {
					drugBankURL, err := queryDrugBankURL(
						ctx,
						e.Query,
						"bing",
						"https://api.bing.microsoft.com/",
						subscriptionKey,
					)
					if err == errFailedToFindDrugBankURL {
						entityLog.Debug("failed to find drugbank url, removing entity")
						if err := querySched.Remove(rawEnt); err != nil {
							return errors.Wrap(err, "failed to remove entity from query scheduler, this will result in multiple repeated requests to youtube")
						}
						return nil
					} else if err != nil {
						return errors.Wrap(err, "failed to query drugbank url")
					}
					entityLog.Debug("found drugbank entry, querying", zap.String("url", drugBankURL))
					start := time.Now()
					drug, err := queryDrug(
						ctx,
						drugBankURL,
					)
					if err != nil {
						return errors.Wrap(err, "failed to query drug")
					}
					entityLog.Debug("queried drugbank", base.Elapsed(start))
					if err := infoCache.SetDrug(
						ctx,
						e.Query,
						drug,
					); err != nil {
						return errors.Wrap(err, "failed to set drug in cache")
					}
					entityLog.Debug("cached drug")
					body, err := json.Marshal(drug)
					if err != nil {
						return err
					}
					if err := pub.Publish(
						ctx,
						drugTopic(e.Query),
						body,
					); err != nil {
						return errors.Wrap(err, "failed to publish drug topic")
					}
				}
				// add the pdb url to the scheduler
				body, err := json.Marshal(&pdbItem{
					Query:                   e.Query,
					DrugBankAccessionNumber: drug.DrugBankAccessionNumber,
					URL:                     drug.Structure.PDB,
					Force:                   e.Force,
				})
				if err != nil {
					return err
				}
				if err := pdbSched.Add(string(body)); err != nil {
					return errors.Wrap(err, "failed to add pdb to scheduler")
				}
				if err := downloadDrugSVG(
					ctx,
					drug.DrugBankAccessionNumber,
					svgCache,
					nil,
				); err != nil {
					return errors.Wrap(err, "downloadDrugSVG")
				}
			default:
				panic(base.Unreachable)
			}
			if err := querySched.Remove(rawEnt); err != nil {
				return errors.Wrap(err, "failed to remove entity from query scheduler, this will result in multiple repeated requests to youtube")
			}
			entityLog.Debug("finished processing drug query")
			return nil
		}(); err != nil {
			entityLog.Warn("query worker error", zap.Error(err))
		}
	}
}
