package transcriber

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"go.uber.org/zap"
)

func readTranscription(
	ctx context.Context,
	events <-chan transcribestreamingservice.MedicalTranscriptResultStreamEvent,
	transcripts chan<- *transcribestreamingservice.MedicalTranscript,
	log *zap.Logger,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev := <-events:
			if ev == nil {
				continue
			}
			if e, ok := ev.(*transcribestreamingservice.MedicalTranscriptEvent); ok {
				if e.Transcript == nil || len(e.Transcript.Results) == 0 {
					continue
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case transcripts <- e.Transcript:
					continue
				default:
					log.Warn("transcript channel full, discarding event")
				}
			} else {
				log.Warn(
					"unrecognized event",
					zap.String("type", fmt.Sprintf("%T", ev)),
				)
			}
		}
	}
}
