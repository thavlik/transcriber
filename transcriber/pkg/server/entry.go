package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/transcriber/pkg/comprehend"

	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	metricsPort int,
	broadcaster base.ServiceOptions,
	specialty string,
	streamKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		broadcaster,
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
			metricsPort,
			log,
		)
	})

	s.spawn(s.pumpNewSource)

	return s.ListenAndServe(
		httpPort,
	)
}
