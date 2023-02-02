package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/transcriber/pkg/refmat"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"github.com/thavlik/transcriber/transcriber/pkg/transcriber"

	"go.uber.org/zap"
)

type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	specialty string
	newSource chan source.Source
	job       *transcriber.TranscriptionJob
	l         chan struct{}
	conns     map[*wsClient]struct{}
	connsL    sync.Mutex
	streamKey string
	refs      refmat.ReferenceMap
	usedRefsL sync.Mutex
	usedRefs  map[*refmat.ReferenceMaterial]time.Time
	wg        *sync.WaitGroup
	log       *zap.Logger
}

func NewServer(
	ctx context.Context,
	specialty string,
	streamKey string,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		specialty: specialty,
		newSource: make(chan source.Source, 16),
		l:         make(chan struct{}, 1),
		conns:     make(map[*wsClient]struct{}),
		streamKey: streamKey,
		refs:      refmat.BuildReferenceMap(refmat.TestReferenceMaterials),
		usedRefs:  make(map[*refmat.ReferenceMaterial]time.Time),
		wg:        new(sync.WaitGroup),
		log:       log,
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
	mux.HandleFunc("/source", s.handleNewSource())
	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("0.0.0.0:%d", httpPort),
		// no read/write timeout because streams must
		// be able to write for a very long time
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
