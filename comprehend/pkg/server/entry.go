package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"

	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	entityCache entitycache.EntityCache,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		entityCache,
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			ctx,
			serverOpts.MetricsPort,
			log,
		)
	})

	base.SignalReady(log)

	return s.ListenAndServe(
		serverOpts.Port,
	)
}
