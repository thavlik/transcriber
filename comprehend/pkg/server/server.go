package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"

	"go.uber.org/zap"
)

type Server struct {
	ctx         context.Context
	cancel      context.CancelFunc
	entityCache entitycache.EntityCache
	wg          *sync.WaitGroup
	log         *zap.Logger
}

func NewServer(
	ctx context.Context,
	entityCache entitycache.EntityCache,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:         ctx,
		cancel:      cancel,
		entityCache: entityCache,
		wg:          new(sync.WaitGroup),
		log:         log,
	}
}

func (s *Server) ListenAndServe(
	httpPort int,
) error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	mux := http.NewServeMux()
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/lookup", s.handleLookup())
	mux.HandleFunc("/detect", s.handleDetect())
	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
		ReadTimeout:  12 * time.Second,
		WriteTimeout: 12 * time.Second,
	}
	s.spawn(func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	})
	s.log.Info(
		"http server listening forever",
		zap.Int("port", httpPort),
	)
	return srv.ListenAndServe()
}

func (s *Server) ShutDown() {
	s.cancel()
	s.wg.Wait()
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}
