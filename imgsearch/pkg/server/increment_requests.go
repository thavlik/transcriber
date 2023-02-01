package server

import (
	"context"

	"go.uber.org/zap"
)

// incrementRequests increments the request counter for the given image
func (s *Server) incrementRequests(metaHash string) {
	s.wg.Add(1)
	defer s.wg.Done()
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()
	if err := s.imageCache.IncrementRequests(
		ctx,
		metaHash,
	); err != nil {
		s.log.Error(
			"failed to increment image request counter",
			zap.Error(err),
			zap.String("hash", metaHash),
		)
	}
}
