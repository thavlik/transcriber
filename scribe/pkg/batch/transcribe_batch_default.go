package batch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

// youtube-dl -o "%(id)s.mp3" -f bestaudio --extract-audio --audio-format mp3 --audio-quality 0 --postprocessor-args '-osr 44100 -ac 1' -- https://www.youtube.com/watch?v=Ipbu796b2_I

// an error string fragment yielded by Amazon Transcribe when a job is not found
const errJobNotFound = "The requested job couldn't be found. Check the job name and try your request again."

func TranscribeBatchDefault(
	ctx context.Context,
	source *BatchTranscribeSource,
	outputBucket string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	svc := transcribeservice.New(base.AWSSession())
	jobName := aws.String(source.Key)
	jobLog := log.With(zap.String("jobName", *jobName))
	var start *time.Time
	if existing, err := svc.GetTranscriptionJobWithContext(
		ctx,
		&transcribeservice.GetTranscriptionJobInput{
			TranscriptionJobName: jobName,
		},
	); (err != nil && strings.Contains(err.Error(), errJobNotFound)) || (err == nil && existing.TranscriptionJob.FailureReason != nil) {
		if err == nil && existing.TranscriptionJob.FailureReason != nil {
			jobLog.Error("transcription job failed, deleting and trying again",
				zap.String("failureReason", *existing.TranscriptionJob.FailureReason))
			if _, err := svc.DeleteTranscriptionJobWithContext(
				ctx,
				&transcribeservice.DeleteTranscriptionJobInput{
					TranscriptionJobName: jobName,
				},
			); err != nil {
				return errors.Wrap(err, "DeleteTranscriptionJobWithContext")
			}
		}
		format := source.Format()
		if format == "" {
			return errors.New("unable to determine source file format")
		}
		resp, err := svc.StartTranscriptionJobWithContext(
			ctx,
			&transcribeservice.StartTranscriptionJobInput{
				TranscriptionJobName: jobName,
				LanguageCode:         aws.String(transcribeservice.LanguageCodeEnUs),
				Media: &transcribeservice.Media{
					MediaFileUri: aws.String(source.S3Uri()),
				},
				OutputBucketName: aws.String(outputBucket),
				OutputKey:        aws.String(strings.Replace(source.Key, "."+format, ".json", 1)),
			})
		if err != nil {
			return errors.Wrap(err, "StartStreamTranscriptionWithContext")
		}
		start = resp.TranscriptionJob.CreationTime
		jobLog.Info("started transcription job", zap.String("format", format))
	} else if err != nil {
		return errors.Wrap(err, "GetTranscriptionJobWithContext")
	} else if existing.TranscriptionJob.CompletionTime != nil {
		fmt.Println("transcription is complete:")
		fmt.Println(*existing.TranscriptionJob.Transcript.TranscriptFileUri)
		return nil
	} else {
		start = existing.TranscriptionJob.CreationTime
		jobLog.Info(
			"transcription job is already running",
			base.Elapsed(*start))
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			existing, err := svc.GetTranscriptionJobWithContext(
				ctx,
				&transcribeservice.GetTranscriptionJobInput{
					TranscriptionJobName: jobName,
				},
			)
			if err != nil {
				return errors.Wrap(err, "GetTranscriptionJobWithContext")
			}
			if existing.TranscriptionJob.FailureReason != nil {
				return errors.Errorf(
					"transcription job failed: %s",
					*existing.TranscriptionJob.FailureReason,
				)
			}
			if existing.TranscriptionJob.CompletionTime != nil {
				fmt.Println("transcription is complete:")
				fmt.Println(*existing.TranscriptionJob.Transcript.TranscriptFileUri)
				return nil
			}
			jobLog.Info(
				"transcription job is still in progress",
				base.Elapsed(*start))
		}
	}
}
