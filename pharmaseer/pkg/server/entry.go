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
	cancelDownloadTopic = "cancel_pdb"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	pubSub pubsub.PubSub,
	querySched scheduler.Scheduler,
	pdbSched scheduler.Scheduler,
	infoCache infocache.InfoCache,
	pdbCache pdbcache.PDBCache,
	concurrency int,
	log *zap.Logger,
) error {
	s := NewServer(
		ctx,
		querySched,
		pdbSched,
		pubSub,
		infoCache,
		pdbCache,
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			s.ctx,
			serverOpts.MetricsPort,
			log,
		)
	})

	stopPopQuery := make(chan struct{}, 1)
	defer func() { stopPopQuery <- struct{}{} }()
	initQueryWorkers(
		concurrency,
		infoCache,
		pubsub.Publisher(pubSub),
		querySched,
		pdbSched,
		stopPopQuery,
		log,
	)

	cancelDownload, err := pubSub.Subscribe(
		s.ctx,
		cancelDownloadTopic,
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to topic")
	}

	stopPopDl := make(chan struct{}, 1)
	defer func() { stopPopDl <- struct{}{} }()
	initDownloadWorkers(
		s.ctx,
		concurrency,
		pdbSched,
		cancelDownload.Messages(s.ctx),
		pdbCache,
		stopPopDl,
		log,
	)
	base.SignalReady(log)
	return s.ListenAndServe(serverOpts.Port)
}
