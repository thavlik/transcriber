package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"github.com/thavlik/transcriber/define/pkg/storage"

	"go.uber.org/zap"
)

type Server struct {
	ctx          context.Context
	cancel       context.CancelFunc
	storage      storage.Storage
	diseaseCache diseasecache.DiseaseCache
	wg           *sync.WaitGroup
	gpt3         gpt3.Client
	log          *zap.Logger
}

func NewServer(
	ctx context.Context,
	storage storage.Storage,
	diseaseCache diseasecache.DiseaseCache,
	gpt3 gpt3.Client,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:          ctx,
		cancel:       cancel,
		storage:      storage,
		diseaseCache: diseaseCache,
		gpt3:         gpt3,
		wg:           new(sync.WaitGroup),
		log:          log,
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
	mux.HandleFunc("/completion", s.handleDefine())
	mux.HandleFunc("/disease", s.handleDisease())

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
	s.wg.Wait()
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}
