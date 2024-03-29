package server

import (
	"context"

	"github.com/thavlik/transcriber/scribe/pkg/transcribe"

	"go.uber.org/zap"
)

func (s *Server) handleTranscript(
	ctx context.Context,
	transcript *transcribe.Transcript,
) error {
	text := transcript.Flatten()
	s.spawn(func() {
		s.broadcastMessage(
			ctx,
			"transcript",
			map[string]interface{}{
				"text": text,
			},
		)
	})
	s.spawn(func() {
		entities, err := s.detectEntities(
			ctx,
			"amazon-comprehend-medical",
			text,
		)
		if err != nil {
			s.log.Error("comprehend error", zap.Error(err))
			return
		}
		if len(entities) == 0 {
			return
		}
		s.broadcastMessage(
			ctx,
			"keyterms",
			map[string]interface{}{
				"entities": entities,
			})
		first := entities[0]
		top := first
		for _, entity := range entities[1:] {
			if entity.Score > top.Score {
				top = entity
			}
		}
		s.log.Debug("comprehended entities",
			zap.Int("count", len(entities)),
			zap.String("first.Text", first.Text),
			zap.Float64("first.Score", first.Score),
			zap.String("top.Text", top.Text),
			zap.Float64("top.Score", top.Score),
		)
	})
	return nil
}
