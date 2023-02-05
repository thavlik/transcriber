package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

type pdbItem struct {
	Query                   string `json:"q"`
	DrugBankAccessionNumber string `json:"d"`
	URL                     string `json:"u"`
	Force                   bool   `json:"f"`
}

func initDownloadWorkers(
	ctx context.Context,
	concurrency int,
	pdbSched scheduler.Scheduler,
	cancelDownload <-chan []byte,
	pdbCache pdbcache.PDBCache,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popPDB := make(chan string)
	go downloadPopper(popPDB, pdbSched, stop, log)
	cancels := make([]chan []byte, concurrency)
	for i := 0; i < concurrency; i++ {
		cancel := make(chan []byte, 8)
		cancels[i] = cancel
		go downloadWorker(
			popPDB,
			cancel,
			pdbSched,
			pdbCache,
			log,
		)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case pdb, ok := <-cancelDownload:
				if !ok {
					return
				}
				for _, cancel := range cancels {
					select {
					case <-ctx.Done():
						return
					case cancel <- pdb:
					}
				}
			}
		}
	}()
}

func downloadPopper(
	popPDB chan<- string,
	pdbSched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := pdbSched.Notify()
	defer close(popPDB)
	delay := 6 * time.Second
	for {
		start := time.Now()
		pdbs, err := pdbSched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(pdbs) > 0 {
			log.Debug("checking pdbs", zap.Int("num", len(pdbs)))
			for _, pdb := range pdbs {
				popPDB <- pdb
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

func downloadWorker(
	popPDB <-chan string,
	cancelDownload <-chan []byte,
	pdbSched scheduler.Scheduler,
	pdbCache pdbcache.PDBCache,
	log *zap.Logger,
) {
	for {
		rawPDB, ok := <-popPDB
		if !ok {
			return
		}
		pdb := new(pdbItem)
		if err := json.Unmarshal([]byte(rawPDB), pdb); err != nil {
			panic(err)
		}
		pdbLog := log.With(
			zap.String("query", pdb.Query),
			zap.String("url", pdb.URL))
		lock, err := pdbSched.Lock(pdb.Query)
		if err == scheduler.ErrLocked {
			// go to the next project
			pdbLog.Debug("pdb already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		pdbLog.Debug("locked pdb")
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			stop := make(chan struct{}, 1)
			stopped := make(chan struct{})
			defer func() {
				cancel()
				_ = lock.Release()
				stop <- struct{}{}
				<-stopped
			}()
			go func() {
				defer func() { stopped <- struct{}{} }()
				for {
					select {
					case cancelPDB := <-cancelDownload:
						if string(cancelPDB) == pdb.Query {
							cancel()
							pdbLog.Debug("download was intentionally cancelled prematurely")
							return
						}
					case <-stop:
						return
					case <-ctx.Done():
						return
					}
				}
			}()
			if err := downloadPDB(
				ctx,
				pdb,
				pdbCache,
				pdbLog,
			); err != nil {
				pdbLog.Error("error downloading pdb", zap.Error(err))
				return
			}
			if err := pdbSched.Remove(rawPDB); err != nil {
				pdbLog.Warn("failed to remove pdb from download scheduler, this will result in multiple repeated requests to drugbank")
			}
		}()
	}
}
