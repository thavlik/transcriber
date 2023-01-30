package transcriber

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
	"github.com/thavlik/transcriber/transcriber/pkg/util"
	"go.uber.org/zap"
)

func Transcribe(
	ctx context.Context,
	source source.Source,
	transcripts chan<- *transcribestreamingservice.MedicalTranscript,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// read stream details
	sampleRate, err := source.SampleRate()
	if err != nil {
		return errors.Wrap(err, "source.SampleRate")
	}
	stereo, err := source.IsStereo()
	if err != nil {
		return errors.Wrap(err, "source.IsStereo")
	}
	log.Debug("starting transcription",
		zap.Int64("sampleRate", sampleRate))
	// start the transcription stream
	var enableChannelIdentification *bool
	if stereo {
		enableChannelIdentification = aws.Bool(true)
	}
	var numberOfChannels *int64
	if stereo {
		numberOfChannels = aws.Int64(2)
	}
	svc := transcribestreamingservice.New(util.AWSSession())
	resp, err := svc.StartMedicalStreamTranscriptionWithContext(
		ctx,
		&transcribestreamingservice.StartMedicalStreamTranscriptionInput{
			LanguageCode:                aws.String("en-US"),
			MediaEncoding:               aws.String("pcm"),
			MediaSampleRateHertz:        aws.Int64(sampleRate),
			NumberOfChannels:            numberOfChannels,
			EnableChannelIdentification: enableChannelIdentification,
			Specialty:                   aws.String("RADIOLOGY"), // PRIMARYCARE | CARDIOLOGY | NEUROLOGY | ONCOLOGY | RADIOLOGY | UROLOGY
			Type:                        aws.String("DICTATION"), // CONVERSATION | DICTATION
		})
	if err != nil {
		return errors.Wrap(err, "StartMedicalStreamTranscriptionWithContext")
	}
	stream := resp.GetStream()

	// spin up a goroutine for sending the audio stream
	writeAudioStreamErr := make(chan error)
	go func() {
		writeAudioStreamErr <- writeAudioStream(
			ctx,
			source,
			stream,
			log,
		)
	}()

	// read the transcription event stream
	readTranscriptionErr := make(chan error)
	go func() {
		readTranscriptionErr <- readTranscription(
			ctx,
			stream.Events(),
			transcripts,
			log,
		)
	}()

	var multi error
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-writeAudioStreamErr:
		cancel()
		multi = multierror.Append(multi, errors.Wrap(err, "writeAudioStream"))
		if err := <-readTranscriptionErr; err != nil {
			multi = multierror.Append(multi, errors.Wrap(err, "readTranscription"))
		}
	case err := <-readTranscriptionErr:
		cancel()
		multi = multierror.Append(multi, errors.Wrap(err, "readTranscription"))
		if err := <-writeAudioStreamErr; err != nil {
			multi = multierror.Append(multi, errors.Wrap(err, "writeAudioStream"))
		}
	}
	return multi
}
