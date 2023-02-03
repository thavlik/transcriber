package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	rtmpPort int,
	scribe base.ServiceOptions,
	streamKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		scribe,
		streamKey,
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

	s.spawn(s.pumpNewSource)

	return s.ListenAndServe(
		serverOpts.Port,
		rtmpPort,
	)
}
