package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func runMetrics(port int, log *zap.Logger) {
	if port == 0 {
		return
	}
	log.Debug("metrics server listening forever", zap.Int("port", port))
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	panic((&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe())
}
