package server

import (
	"go.uber.org/zap"
)

// incrementRequests increments the request counter for the given image
func (s *Server) incrementRequests(metaHash string) {
	if err := s.imageCache.IncrementRequests(
		s.ctx,
		metaHash,
	); err != nil {
		s.log.Error(
			"failed to increment image request counter",
			zap.Error(err),
			zap.String("hash", metaHash),
		)
	}
}
