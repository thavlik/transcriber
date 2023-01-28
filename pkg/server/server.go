package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thavlik/transcriber/pkg/refmat"
	"github.com/thavlik/transcriber/pkg/source"
	"github.com/thavlik/transcriber/pkg/transcriber"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type server struct {
	newSource chan source.Source
	job       *transcriber.TranscriptionJob
	l         chan struct{}
	conns     map[*websocket.Conn]struct{}
	connsL    sync.Mutex
	streamKey string
	refs      refmat.ReferenceMap
	usedRefsL sync.Mutex
	usedRefs  map[*refmat.ReferenceMaterial]struct{}
	log       *zap.Logger
}

func (s *server) ListenAndServe(
	httpPort int,
	rtmpPort int,
) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/ws", s.handleWebSock())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		s.log.Warn("404", zap.String("r.RequestURI", r.RequestURI))
	})
	httpDone := make(chan error)
	go func() {
		s.log.Info("http server listening forever", zap.Int("port", httpPort))
		httpDone <- (&http.Server{
			Handler:      mux,
			Addr:         fmt.Sprintf("0.0.0.0:%d", httpPort),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}).ListenAndServe()
	}()
	rtmpDone := make(chan error)
	go func() {
		rtmpDone <- s.listenRTMP(rtmpPort)
	}()
	select {
	case err := <-httpDone:
		return errors.Wrap(err, "http server failed")
	case err := <-rtmpDone:
		return errors.Wrap(err, "rtmp server failed")
	}
}
