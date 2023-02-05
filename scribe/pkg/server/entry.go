package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"

	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	broadcasterOpts *base.ServiceOptions,
	comprehendOpts *base.ServiceOptions,
	specialty string,
	streamKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		broadcasterOpts,
		comprehendOpts,
		specialty,
		streamKey,
		&comprehend.Filter{
			// TODO: make exclude terms configurable
			ExcludeTerms: []string{"Um", "Uh", "Uhm"},
		},
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

	base.SignalReady(log)

	return s.ListenAndServe(
		serverOpts.Port,
	)
}
