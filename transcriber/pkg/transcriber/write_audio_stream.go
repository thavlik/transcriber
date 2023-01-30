package transcriber

import (
	"context"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"github.com/thavlik/transcriber/transcriber/pkg/util"
	"go.uber.org/zap"
)

func writeAudioStream(
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
			n, err := source.ReadAudioChunk(buf)
			if err != nil {
				return errors.Wrap(err, "source.ReadAudioChunk")
			} else if n == 0 {
				// no audio data available
				log.Warn("source returned no audio data")
				continue
			}
			//log.Debug("read pcm audio chunk", zap.Int("bytes", n))
			if err := stream.Send(
				ctx,
				&transcribestreamingservice.AudioEvent{
					// it's wise to duplicate the buffer in case
					// the memory is reused by the source before
					// the stream has finished sending it
					AudioChunk: util.Duplicate(buf[:n]),
				},
			); err != nil {
				return errors.Wrap(err, "stream.Send")
			}
			//log.Debug("sent pcm audio chunk", zap.Int("bytes", n))
		}
	}
}