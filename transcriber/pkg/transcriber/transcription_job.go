package transcriber

import (
	"context"

	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"go.uber.org/zap"
)

type TranscriptionJob struct {
	ctx         context.Context
	cancel      context.CancelFunc
	source      source.Source
	specialty   string
	transcripts chan<- *transcribestreamingservice.MedicalTranscript
	log         *zap.Logger
}

func NewTranscriptionJob(
	ctx context.Context,
	source source.Source,
	specialty string,
	transcripts chan<- *transcribestreamingservice.MedicalTranscript,
	log *zap.Logger,
) *TranscriptionJob {
	childCtx, cancel := context.WithCancel(ctx)
	return &TranscriptionJob{
		ctx:         childCtx,
		cancel:      cancel,
		source:      source,
		specialty:   specialty,
		transcripts: transcripts,
		log:         log,
	}
}

func (j *TranscriptionJob) Context() context.Context {
	return j.ctx
}

func (j *TranscriptionJob) Cancel() {
	j.cancel()
}

func (j *TranscriptionJob) Transcribe() error {
	return Transcribe(
		j.ctx,
		j.source,
		j.specialty,
		j.transcripts,
		j.log,
	)
}
