package base

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func RunMetrics(
	ctx context.Context,
	port int,
	log *zap.Logger,
) {
	if port == 0 {
		return
	}
	log.Debug(
		"metrics server listening forever",
		zap.Int("port", port),
	)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	}()
	panic(srv.ListenAndServe())
}
