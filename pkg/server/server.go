package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/yutopp/go-rtmp"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Entry(
	httpPort int,
	rtmpPort int,
	metricsPort int,
	log *zap.Logger,
) error {
	go runMetrics(metricsPort, log)
	return (&server{
		log,
	}).ListenAndServe(
		httpPort,
		rtmpPort,
	)
}

type server struct {
	log *zap.Logger
}

func (s *server) ListenAndServe(
	httpPort int,
	rtmpPort int,
) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {})
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

func (s *server) listenRTMP(port int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "resolve failed")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}
	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			return conn, &rtmp.ConnConfig{
				Handler: &Handler{log: s.log},
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
			}
		},
	})
	s.log.Info("rtmp server listening forever", zap.Int("port", port))
	if err := srv.Serve(listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}
