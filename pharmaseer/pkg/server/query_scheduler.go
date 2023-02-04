package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"go.uber.org/zap"
)

type entityType int

func (e entityType) String() string {
	switch e {
	case drug:
		return "drug"
	default:
		panic(base.Unreachable)
	}
}

const (
	drug entityType = 1
)

type entity struct {
	Type entityType `json:"t"`
	ID   string     `json:"_"`
}

func initQueryWorkers(
	concurrency int,
	infoCache infocache.InfoCache,
	pub pubsub.Publisher,
	querySched scheduler.Scheduler,
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
	pub pubsub.Publisher,
	log *zap.Logger,
) {
	for {
		rawEnt, ok := <-popQuery
		if !ok {
			return
		}
		e := &entity{}
		if err := json.Unmarshal([]byte(rawEnt), e); err != nil {
			fmt.Println("removing malformed rawEnt")
			fmt.Println(string(rawEnt))
			if err := querySched.Remove(rawEnt); err != nil {
				panic(err)
			}
			continue
		}
		entityLog := log.With(
			zap.String("type", e.Type.String()),
			zap.String("id", e.ID))
		lock, err := querySched.Lock(e.ID)
		if err == scheduler.ErrLocked {
			// go to the next project
			entityLog.Debug("entity already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		entityLog.Debug("processing query item")
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
			case drug:
			default:
				panic(base.Unreachable)
			}
			if err := querySched.Remove(rawEnt); err != nil {
				return errors.Wrap(err, "failed to remove entity from query scheduler, this will result in multiple repeated requests to youtube")
			}
			return nil
		}(); err != nil {
			entityLog.Warn("query worker error", zap.Error(err))
		}
	}
}
