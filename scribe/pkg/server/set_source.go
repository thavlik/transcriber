package server

import (
	"context"

	"github.com/thavlik/transcriber/scribe/pkg/source"
	"github.com/thavlik/transcriber/scribe/pkg/transcribe"

	"go.uber.org/zap"
)

func (s *Server) setSource(
	ctx context.Context,
	src source.Source,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.l <- struct{}{}:
		defer func() { <-s.l }()
		if s.job != nil {
			// kill the old job
			s.job.Cancel()
		}
		// start a new job, use the source's context
		// so the transcription is cancelled when the
		// source is closed
		transcripts := make(chan *transcribe.Transcript, 16)
		s.job = transcribe.NewTranscriptionJob(
			src.Context(),
			src,
			s.specialty,
			transcripts,
			s.log,
		)
		s.spawn(func() {
			for {
				select {
				case <-s.job.Context().Done():
					return
				case transcript := <-transcripts:
					if err := s.handleTranscript(
						ctx,
						transcript,
					); err != nil {
						s.log.Error("handle transcript error", zap.Error(err))
					}
				}
			}
		})
		s.spawn(func() {
			s.log.Debug("starting transcribe job")
			if err := s.job.Transcribe(); err != nil {
				s.log.Error("transcribe error", zap.Error(err))
			}
		})
		return nil
	}
}
