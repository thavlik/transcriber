package server

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type server struct {
	apiKey   string
	endpoint string
	log      *zap.Logger
}

func newServer(
	apiKey string,
	endpoint string,
	log *zap.Logger,
) *server {
	return &server{
		apiKey:   apiKey,
		endpoint: endpoint,
		log:      log,
	}
}

func (s *server) listenAndServe(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/", s.handleSearch())
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
