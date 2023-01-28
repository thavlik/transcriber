package server

import (
	"context"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/thavlik/transcriber/pkg/source"
	"github.com/thavlik/transcriber/pkg/transcriber"

	"go.uber.org/zap"
)

func (s *server) setSource(
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
		transcripts := make(chan *transcribestreamingservice.MedicalTranscript, 16)
		s.job = transcriber.NewTranscriptionJob(
			src.Context(),
			src,
			transcripts,
			s.log,
		)
		go func() {
			for {
				select {
				case <-s.job.Context().Done():
					return
				case transcript := <-transcripts:
					if err := s.handleTranscript(transcript); err != nil {
						s.log.Error("handle transcript error", zap.Error(err))
					}
				}
			}
		}()
		go func() {
			s.log.Debug("starting transcribe job")
			if err := s.job.Transcribe(); err != nil {
				s.log.Error("transcribe error", zap.Error(err))
			}
		}()
		return nil
	}
}
