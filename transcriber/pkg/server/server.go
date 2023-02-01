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

	"github.com/pkg/errors"
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
	rtmpPort int,
) error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", base.Handle200)
	mux.HandleFunc("/readyz", base.Handle200)
	mux.HandleFunc("/ws", s.handleWebSock())

	srv := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	httpDone := make(chan error, 1)
	s.spawn(func() {
		s.log.Info(
			"http server listening forever",
			zap.Int("port", httpPort),
		)
		httpDone <- srv.ListenAndServe()
	})

	s.spawn(func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	})

	rtmpDone := make(chan error, 1)
	s.spawn(func() {
		rtmpDone <- s.listenRTMP(ctx, rtmpPort)
	})

	select {
	case err := <-httpDone:
		cancel()
		return errors.Wrap(err, "http server failed")
	case err := <-rtmpDone:
		cancel()
		return errors.Wrap(err, "rtmp server failed")
	}
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
