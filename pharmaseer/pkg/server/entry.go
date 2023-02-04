package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

var (
	cancelVideoTopic = "cancel_video"
)

func Entry(
	ctx context.Context,
	port int,
	pubSub pubsub.PubSub,
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	infoCache infocache.InfoCache,
	pdbCache pdbcache.PDBCache,
	concurrency int,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		querySched,
		dlSched,
		pubSub,
		infoCache,
		pdbCache,
		log,
	)

	stopPopQuery := make(chan struct{}, 1)
	defer func() { stopPopQuery <- struct{}{} }()
	initQueryWorkers(
		concurrency,
		infoCache,
		pubsub.Publisher(pubSub),
		querySched,
		stopPopQuery,
		log,
	)

	cancelVideoDownload, err := pubSub.Subscribe(
		ctx,
		cancelVideoTopic,
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to topic")
	}

	stopPopDl := make(chan struct{}, 1)
	defer func() { stopPopDl <- struct{}{} }()
	initDownloadWorkers(
		ctx,
		concurrency,
		dlSched,
		cancelVideoDownload.Messages(ctx),
		pdbCache,
		stopPopDl,
		log,
	)

	base.SignalReady(log)
	return s.ListenAndServe(port)
}
