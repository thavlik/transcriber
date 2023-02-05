package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"github.com/thavlik/transcriber/pharmaseer/pkg/thumbcache"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	querySched scheduler.Scheduler
	dlSched    scheduler.Scheduler
	pubsub     pubsub.PubSub
	infoCache  infocache.InfoCache
	pdbCache   pdbcache.PDBCache
	svgCache   thumbcache.ThumbCache
	log        *zap.Logger
}

func NewServer(
	ctx context.Context,
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	pub pubsub.PubSub,
	infoCache infocache.InfoCache,
	pdbCache pdbcache.PDBCache,
	svgCache thumbcache.ThumbCache,
	log *zap.Logger,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx,
		cancel,
		new(sync.WaitGroup),
		querySched,
		dlSched,
		pub,
		infoCache,
		pdbCache,
		svgCache,
		log,
	}
}

func (s *Server) spawn(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}

func (s *Server) ShutDown() {
	s.cancel()
	s.wg.Wait()
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterPharmaSeer(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/structure", s.handleStructure())
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}).ListenAndServe()
}
