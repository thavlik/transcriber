package server

import (
	"context"
	"time"

	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"go.uber.org/zap"
)

// a goroutine for assigning a new source
func (s *Server) pumpNewSource() {
	var ctx context.Context
	var cancel context.CancelFunc
	for {
		select {
		case <-s.ctx.Done():
			if cancel != nil {
				cancel()
			}
			return
		case newSource := <-s.newSource:
			if cancel != nil {
				cancel()
			}
			ctx, cancel = context.WithCancel(newSource.Context())
			go func(ctx context.Context, src source.Source) {
				retryDelay := time.Second
				for {
					s.log.Debug("pushing audio source")
					if err := s.pushAudioSource(
						ctx,
						src,
					); err != nil {
						s.log.Error(
							"failed to push audio source",
							zap.Error(err),
							zap.String("retryDelay", retryDelay.String()),
						)
						select {
						case <-ctx.Done():
							return
						case <-time.After(retryDelay):
							continue
						}
					}
				}
			}(ctx, newSource)
			s.log.Info("received new audio source")
		}
	}
}
