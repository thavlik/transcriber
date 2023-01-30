package base

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func RunMetrics(port int, log *zap.Logger) {
	if port == 0 {
		return
	}
	log.Debug("metrics server listening forever", zap.Int("port", port))
	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
