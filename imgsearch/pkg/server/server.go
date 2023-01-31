package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"go.uber.org/zap"
)

type server struct {
	apiKey     string
	endpoint   string
	imageCache *cache.ImageCache
	log        *zap.Logger
}

func newServer(
	apiKey string,
	endpoint string,
	imageCache *cache.ImageCache,
	log *zap.Logger,
) *server {
	return &server{
		apiKey:     apiKey,
		endpoint:   endpoint,
		imageCache: imageCache,
		log:        log,
	}
}

func (s *server) listenAndServe(port int) error {
	mux := http.NewServeMux()
	// it's okay if the health and ready checks are publicly accessible
	mux.HandleFunc("/", base.Handle404(s.log))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/api/search", s.handleSearch())
	mux.HandleFunc("/api/img", s.handleImage())
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
