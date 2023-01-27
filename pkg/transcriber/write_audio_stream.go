package transcriber

import (
	"context"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pkg/source"
)

func writeAudioStream(
	ctx context.Context,
	source source.Source,
	stream *transcribestreamingservice.StartStreamTranscriptionEventStream,
) error {
	buf := make([]byte, 32000) // Amazon has a 32kb max
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := source.ReadAudioChunk(buf)
			if err != nil {
				return errors.Wrap(err, "source.ReadAudioChunk")
			}
			if err := stream.Send(
				ctx,
				&transcribestreamingservice.AudioEvent{
					// it's wise to duplicate the buffer in case
					// the memory is reused by the source before
					// the stream has finished sending it
					AudioChunk: duplicate(buf[:n]),
				},
			); err != nil {
				return errors.Wrap(err, "stream.Send")
			}
		}
	}
}
