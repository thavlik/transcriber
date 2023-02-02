package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	rtmpPort int,
	metricsPort int,
	transcriber base.ServiceOptions,
	streamKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		transcriber,
		streamKey,
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			ctx,
			metricsPort,
			log,
		)
	})

	s.spawn(s.pumpNewSource)

	return s.ListenAndServe(
		httpPort,
		rtmpPort,
	)
}
