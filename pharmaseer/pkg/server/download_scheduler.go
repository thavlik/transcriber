package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

func initDownloadWorkers(
	ctx context.Context,
	concurrency int,
	dlSched scheduler.Scheduler,
	cancelVideoDownload <-chan []byte,
	pdbCache pdbcache.PDBCache,
	stop <-chan struct{},
	log *zap.Logger,
) {
	popVideoID := make(chan string)
	go downloadPopper(popVideoID, dlSched, stop, log)
	cancels := make([]chan []byte, concurrency)
	for i := 0; i < concurrency; i++ {
		cancel := make(chan []byte, 8)
		cancels[i] = cancel
		go downloadWorker(
			popVideoID,
			cancel,
			dlSched,
			pdbCache,
			log,
		)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case videoID, ok := <-cancelVideoDownload:
				if !ok {
					return
				}
				for _, cancel := range cancels {
					select {
					case <-ctx.Done():
						return
					case cancel <- videoID:
					}
				}
			}
		}
	}()
}

func downloadPopper(
	popVideoID chan<- string,
	dlSched scheduler.Scheduler,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := dlSched.Notify()
	defer close(popVideoID)
	delay := 12 * time.Second
	for {
		start := time.Now()
		videoIDs, err := dlSched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(videoIDs) > 0 {
			log.Debug("checking videos", zap.Int("num", len(videoIDs)))
			for _, videoID := range videoIDs {
				popVideoID <- videoID
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
	popVideoID <-chan string,
	cancelVideoDownload <-chan []byte,
	dlSched scheduler.Scheduler,
	pdbCache pdbcache.PDBCache,
	log *zap.Logger,
) {
	for {
		videoID, ok := <-popVideoID
		if !ok {
			return
		}
		videoLog := log.With(zap.String("videoID", videoID))
		lock, err := dlSched.Lock(videoID)
		if err == scheduler.ErrLocked {
			// go to the next project
			videoLog.Debug("video already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		videoLog.Debug("locked video")
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			stop := make(chan struct{}, 1)
			stopped := make(chan struct{})
			defer func() {
				stop <- struct{}{}
				<-stopped
				_ = lock.Release()
				cancel()
			}()
			onProgress := make(chan *base.DownloadProgress, 1)
			defer close(onProgress)
			go func() {
				defer func() { stopped <- struct{}{} }()
				done := ctx.Done()
				for {
					select {
					case cancelID := <-cancelVideoDownload:
						if string(cancelID) == videoID {
							cancel()
							videoLog.Debug("download was intentionally cancelled prematurely")
							return
						}
					case <-stop:
						return
					case <-done:
						return
					case progress, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							videoLog.Warn("failed to extend video download lock", zap.Error(err))
							cancel()
							return
						}
						if progress != nil {
							// TODO
						}
					}
				}
			}()
			base.ProgressDownload(ctx, onProgress)
		}()
	}
}
