package server

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	port int,
	metricsPort int,
	history history.History,
	apiKey string,
	endpoint string,
	imageCache *cache.ImageCache,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		history,
		apiKey,
		endpoint,
		imageCache,
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

	return s.ListenAndServe(port)
}
