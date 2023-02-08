package transcribe

import (
	"context"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/scribe/pkg/source"
	"go.uber.org/zap"
)

func writeAudioStreamMedical(
	ctx context.Context,
	source source.Source,
	stream *transcribestreamingservice.StartMedicalStreamTranscriptionEventStream,
	log *zap.Logger,
) error {
	buf := make([]byte, 32000) // Amazon has a 32kb max
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := source.Read(buf)
			if err != nil {
				return errors.Wrap(err, "source.ReadAudioChunk")
			} else if n == 0 {
				// no audio data available
				log.Warn("source returned no audio data")
				continue
			}
			if err := stream.Send(
				ctx,
				&transcribestreamingservice.AudioEvent{
					// it's wise to duplicate the buffer in case
					// the memory is reused by the source before
					// the stream has finished sending it
					AudioChunk: base.Duplicate(buf[:n]),
				},
			); err != nil {
				return errors.Wrap(err, "stream.Send")
			}
		}
	}
}

func writeAudioStream(
	ctx context.Context,
	source source.Source,
	stream *transcribestreamingservice.StartStreamTranscriptionEventStream,
	log *zap.Logger,
) error {
	buf := make([]byte, 32000) // Amazon has a 32kb max
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := source.Read(buf)
			if err != nil {
				return errors.Wrap(err, "source.ReadAudioChunk")
			} else if n == 0 {
				// no audio data available
				log.Warn("source returned no audio data")
				continue
			}
			if err := stream.Send(
				ctx,
				&transcribestreamingservice.AudioEvent{
					// it's wise to duplicate the buffer in case
					// the memory is reused by the source before
					// the stream has finished sending it
					AudioChunk: base.Duplicate(buf[:n]),
				},
			); err != nil {
				return errors.Wrap(err, "stream.Send")
			}
		}
	}
}
