package server

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/thavlik/transcriber/pkg/source"

	"go.uber.org/zap"
)

func Entry(
	httpPort int,
	rtmpPort int,
	metricsPort int,
	streamKey string,
	log *zap.Logger,
) error {
	go runMetrics(metricsPort, log)
	s := &server{
		newSource: make(chan source.Source, 16),
		l:         make(chan struct{}, 1),
		conns:     make(map[*websocket.Conn]struct{}),
		streamKey: streamKey,
		log:       log,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case src := <-s.newSource:
				log.Info("new source")
				if err := s.setSource(
					ctx,
					src,
				); err != nil {
					return
				}
			}
		}
	}()
	return s.ListenAndServe(
		httpPort,
		rtmpPort,
	)
}
