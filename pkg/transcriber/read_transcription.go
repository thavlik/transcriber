package transcriber

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"go.uber.org/zap"
)

func readTranscription(
	ctx context.Context,
	events <-chan transcribestreamingservice.TranscriptResultStreamEvent,
	transcripts chan<- *transcribestreamingservice.Transcript,
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
			if e, ok := ev.(*transcribestreamingservice.TranscriptEvent); ok {
				log.Debug(
					"transcript event",
					zap.String("text", string(ConvertTranscript(e.Transcript))))
				//select {
				//case <-ctx.Done():
				//	return ctx.Err()
				//case transcripts <- e.Transcript:
				//	continue
				//default:
				//	log.Warn("transcript channel full, discarding event")
				//}
			} else {
				log.Warn(
					"unrecognized event",
					zap.String("type", fmt.Sprintf("%T", ev)),
				)
			}
		}
	}
}
