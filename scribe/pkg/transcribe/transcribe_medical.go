package transcribe

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/scribe/pkg/source"
	"go.uber.org/zap"
)

func TranscribeMedical(
	ctx context.Context,
	source source.Source,
	specialty string,
	transcripts chan<- *Transcript,
	log *zap.Logger,
) error {
	if specialty == "" {
		return errors.New("specialty is required")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sampleRate, err := source.SampleRate()
	if err != nil {
		return errors.Wrap(err, "source.SampleRate")
	} else if sampleRate < minSampleRate {
		return errors.Errorf(
			"sample rate of %d Hz is too low (minimum required is %d Hz)",
			sampleRate,
			minSampleRate,
		)
	}

	isStereo, err := source.IsStereo()
	if err != nil {
		return errors.Wrap(err, "source.IsStereo")
	}

	log.Debug("starting transcription",
		zap.Int64("sampleRate", sampleRate),
		zap.Bool("isStero", isStereo))
	var enableChannelIdentification *bool
	if isStereo {
		// the only acceptable values are nil and true
		enableChannelIdentification = aws.Bool(true)
	}
	var numberOfChannels *int64
	if isStereo {
		// the only acceptable values are nil and 2
		numberOfChannels = aws.Int64(2)
	}
	svc := transcribestreamingservice.New(base.AWSSession())
	resp, err := svc.StartMedicalStreamTranscriptionWithContext(
		ctx,
		&transcribestreamingservice.StartMedicalStreamTranscriptionInput{
			LanguageCode:                aws.String("en-US"),
			MediaEncoding:               aws.String("pcm"),
			MediaSampleRateHertz:        aws.Int64(sampleRate),
			NumberOfChannels:            numberOfChannels,
			EnableChannelIdentification: enableChannelIdentification,
			Specialty:                   aws.String(specialty),   // PRIMARYCARE | CARDIOLOGY | NEUROLOGY | ONCOLOGY | RADIOLOGY | UROLOGY
			Type:                        aws.String("DICTATION"), // CONVERSATION | DICTATION
		})
	if err != nil {
		return errors.Wrap(err, "StartMedicalStreamTranscriptionWithContext")
	}
	stream := resp.GetStream()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	defer wg.Wait()

	// spin up a goroutine for sending the audio stream
	writeAudioStreamErr := make(chan error)
	go func() {
		defer wg.Done()
		writeAudioStreamErr <- writeAudioStreamMedical(
			ctx,
			source,
			stream,
			log,
		)
	}()

	// read the transcription event stream
	readTranscriptionErr := make(chan error)
	go func() {
		defer wg.Done()
		readTranscriptionErr <- readTranscriptionMedical(
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
		multi = multierror.Append(multi, errors.Wrap(err, "writeAudioStreamMedical"))
		if err := <-readTranscriptionErr; err != nil {
			multi = multierror.Append(multi, errors.Wrap(err, "readTranscriptionMedical"))
		}
	case err := <-readTranscriptionErr:
		cancel()
		multi = multierror.Append(multi, errors.Wrap(err, "readTranscriptionMedical"))
		if err := <-writeAudioStreamErr; err != nil {
			multi = multierror.Append(multi, errors.Wrap(err, "writeAudioStreamMedical"))
		}
	}
	return multi
}
