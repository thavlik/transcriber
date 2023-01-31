package server

import (
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"go.uber.org/zap"
)

func Entry(
	port int,
	metricsPort int,
	apiKey string,
	endpoint string,
	imageCache *cache.ImageCache,
	log *zap.Logger,
) error {
	go base.RunMetrics(metricsPort, log)
	return newServer(
		apiKey,
		endpoint,
		imageCache,
		log,
	).listenAndServe(port)
}
