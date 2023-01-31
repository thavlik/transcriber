package server

import (
	"context"

	"go.uber.org/zap"
)

func (s *server) incrementRequests(imageHash string) {
	// Increment the request counter for this image
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.imageCache.IncrementRequests(
		ctx,
		imageHash,
	); err != nil {
		s.log.Error(
			"failed to increment image request counter",
			zap.Error(err),
			zap.String("hash", imageHash),
		)
	}
}
