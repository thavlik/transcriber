package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"go.uber.org/zap"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	sched      scheduler.Scheduler
	pharmaSeer *base.ServiceOptions
	wg         *sync.WaitGroup
	log        *zap.Logger
}

func NewServer(
	ctx context.Context,
	sched scheduler.Scheduler,
	pharmaSeer *base.ServiceOptions,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:        ctx,
		cancel:     cancel,
		sched:      sched,
		pharmaSeer: pharmaSeer,
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
	mux.HandleFunc("/convert", s.handleConvert())
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
