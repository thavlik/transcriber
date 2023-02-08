package transcribe

import (
	"context"

	"github.com/thavlik/transcriber/scribe/pkg/source"
	"go.uber.org/zap"
)

const minSampleRate = int64(8000)

func Transcribe(
	ctx context.Context,
	source source.Source,
	specialty string,
	transcripts chan<- *Transcript,
	log *zap.Logger,
) error {
	if specialty != "" {
		return TranscribeMedical(
			ctx,
			source,
			specialty,
			transcripts,
			log,
		)
	}
	return TranscribeDefault(
		ctx,
		source,
		transcripts,
		log,
	)
}
