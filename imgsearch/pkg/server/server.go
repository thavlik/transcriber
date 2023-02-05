package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	"go.uber.org/zap"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	history    history.History
	apiKey     string
	endpoint   string
	imageCache *cache.ImageCache
	wg         *sync.WaitGroup
	log        *zap.Logger
}

func NewServer(
	ctx context.Context,
	history history.History,
	apiKey string,
	endpoint string,
	imageCache *cache.ImageCache,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:        ctx,
		cancel:     cancel,
		history:    history,
		apiKey:     apiKey,
		endpoint:   endpoint,
		imageCache: imageCache,
		wg:         new(sync.WaitGroup),
		log:        log,
	}
}

func (s *Server) ShutDown() {
	s.cancel()
	s.wg.Wait()
}

func (s *Server) ListenAndServe(port int) error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	mux := http.NewServeMux()
	// it's okay if the health and ready checks are publicly accessible
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/img/search", s.handleSearch())
	mux.HandleFunc("/img/view", s.handleImage())
	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.spawn(func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	})
	s.log.Info("listening forever", zap.Int("port", port))
	return srv.ListenAndServe()
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}
