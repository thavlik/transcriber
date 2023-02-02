package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"

	"go.uber.org/zap"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	conns  map[*wsClient]struct{}
	connsL sync.Mutex
	pub    pubsub.Publisher
	wg     *sync.WaitGroup
	log    *zap.Logger
}

func NewServer(
	ctx context.Context,
	pub pubsub.Publisher,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:    ctx,
		cancel: cancel,
		pub:    pub,
		wg:     new(sync.WaitGroup),
		conns:  make(map[*wsClient]struct{}),
		log:    log,
	}
}

func (s *Server) ListenAndServe(
	httpPort int,
) error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.Handle200)
	mux.HandleFunc("/readyz", base.Handle200)
	mux.HandleFunc("/ws", s.handleWebSock())
	mux.HandleFunc("/publish", s.handlePublish())

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
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
	s.connsL.Lock()
	for conn := range s.conns {
		_ = conn.c.Close()
	}
	s.connsL.Unlock()
	s.wg.Wait()
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}
