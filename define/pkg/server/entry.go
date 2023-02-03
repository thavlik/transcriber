package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	metricsPort int,
	openAISecretKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		openAISecretKey,
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

	return s.ListenAndServe(
		httpPort,
	)
}
