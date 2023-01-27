package transcriber

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pkg/source"
)

func Transcribe(
	ctx context.Context,
	source source.Source,
) error {
	// read stream details
	sampleRate, err := source.SampleRate()
	if err != nil {
		return errors.Wrap(err, "source.SampleRate")
	}
	encoding, err := source.Encoding()
	if err != nil {
		return errors.Wrap(err, "source.Encoding")
	}
	fmt.Printf("sampleRate: %d\n", sampleRate)
	fmt.Printf("encoding: %s\n", encoding)

	// start the transcription stream
	svc := transcribestreamingservice.New(AWSSession())
	resp, err := svc.StartStreamTranscriptionWithContext(
		ctx,
		&transcribestreamingservice.StartStreamTranscriptionInput{
			LanguageCode:         aws.String("en-US"),
			MediaEncoding:        aws.String(encoding),
			MediaSampleRateHertz: aws.Int64(sampleRate),
		})
	if err != nil {
		return errors.Wrap(err, "StartStreamTranscriptionWithContext")
	}
	stream := resp.GetStream()

	// spin up a goroutine for sending the audio stream
	stopped := make(chan error)
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		stopped <- writeAudioStream(
			childCtx,
			source,
			stream,
		)
	}()

	// read the transcription event stream
	if err := readTranscription(ctx, stream.Events()); err != nil {
		return errors.Wrap(err, "readTranscription")
	}

	// wait for the audio stream to finish
	cancel()
	if err := <-stopped; err != nil {
		return errors.Wrap(err, "writeAudioStream")
	}

	return nil
}
