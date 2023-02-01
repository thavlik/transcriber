package server

import (
	"go.uber.org/zap"
)

// a goroutine for assigning a new source
func (s *Server) handleNewSource() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case src := <-s.newSource:
			s.log.Info("received new audio source")
			s.clearUsedRefs()
			if err := s.setSource(
				s.ctx,
				src,
			); err != nil {
				s.log.Error("failed to set audio source", zap.Error(err))
				continue
			}
		}
	}
}
