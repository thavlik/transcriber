package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	"go.uber.org/zap"
)

func (s *Server) pushSearchHistory(query, userID string) {
	if err := s.history.Push(
		s.ctx,
		&history.Search{
			ID:        uuid.New().String(),
			Query:     query,
			UserID:    userID,
			Timestamp: time.Now(),
		},
	); err != nil {
		s.log.Error(
			"failed to push search history",
			zap.String("userID", userID),
			zap.Error(err),
		)
	}
}
