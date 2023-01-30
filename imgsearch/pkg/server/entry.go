package server

import (
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func Entry(
	port int,
	metricsPort int,
	apiKey string,
	log *zap.Logger,
) error {
	go base.RunMetrics(metricsPort, log)
	return newServer(
		apiKey,
		log,
	).listenAndServe(port)
}
