package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/storage"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	metricsPort int,
	storage storage.Storage,
	openAISecretKey string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		storage,
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
