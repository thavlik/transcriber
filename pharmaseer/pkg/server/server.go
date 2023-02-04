package server

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
)

type Server struct {
	querySched scheduler.Scheduler
	dlSched    scheduler.Scheduler
	pubsub     pubsub.PubSub
	infoCache  infocache.InfoCache
	pdbCache   pdbcache.PDBCache
	log        *zap.Logger
}

func NewServer(
	querySched scheduler.Scheduler,
	dlSched scheduler.Scheduler,
	pub pubsub.PubSub,
	infoCache infocache.InfoCache,
	pdbCache pdbcache.PDBCache,
	log *zap.Logger,
) *Server {
	return &Server{
		querySched,
		dlSched,
		pub,
		infoCache,
		pdbCache,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterPharmaSeer(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/healthz", base.HealthHandler)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Hour,
	}).ListenAndServe()
}
