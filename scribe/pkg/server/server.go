package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"github.com/thavlik/transcriber/scribe/pkg/source"
	"github.com/thavlik/transcriber/scribe/pkg/transcribe"

	"go.uber.org/zap"
)

type Server struct {
	ctx         context.Context
	cancel      context.CancelFunc
	broadcaster *base.ServiceOptions
	comprehend  *base.ServiceOptions
	specialty   string
	newSource   chan source.Source
	job         *transcribe.TranscriptionJob
	l           chan struct{}
	filter      *comprehend.Filter
	streamKey   string
	wg          *sync.WaitGroup
	log         *zap.Logger
}

func NewServer(
	ctx context.Context,
	broadcasterOpts *base.ServiceOptions,
	comprehendOpts *base.ServiceOptions,
	specialty string,
	streamKey string,
	filter *comprehend.Filter,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:         ctx,
		cancel:      cancel,
		broadcaster: broadcasterOpts,
		comprehend:  comprehendOpts,
		specialty:   specialty,
		newSource:   make(chan source.Source, 16),
		l:           make(chan struct{}, 1),
		streamKey:   streamKey,
		wg:          new(sync.WaitGroup),
		filter:      filter,
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
	mux.HandleFunc("/healthz", base.Handle200)
	mux.HandleFunc("/readyz", base.Handle200)
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
