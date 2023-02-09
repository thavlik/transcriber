package batch

import (
	"context"
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

type Config struct {
	// Source is the source of the audio to transcribe.
	// It must be an object in S3, on AWS specifically.
	Source *BatchTranscribeSource `json:"source" yaml:"source"`

	// OutputBucket is the S3 bucket to write the transcription output to.
	OutputBucket string `json:"outputBucket" yaml:"outputBucket"`

	// MaxSpeakerLabels is the maximum number of speakers to detect in the audio.
	// 0 = no speaker labels
	MaxSpeakerLabels int64 `json:"maxSpeakerLabels" yaml:"maxSpeakerLabels"`
}

func TranscribeBatchDefault(
	ctx context.Context,
	config *Config,
	follow bool,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	svc := transcribeservice.New(base.AWSSession())
	jobName := aws.String(config.OutputBucket + "-" + config.Source.Key)
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
		format := config.Source.Format()
		if format == "" {
			return errors.New("unable to determine source file format")
		}
		settings := new(transcribeservice.Settings)
		if config.MaxSpeakerLabels > 0 {
			settings.ShowSpeakerLabels = aws.Bool(true)
			settings.MaxSpeakerLabels = aws.Int64(int64(config.MaxSpeakerLabels))
		}
		resp, err := svc.StartTranscriptionJobWithContext(
			ctx,
			&transcribeservice.StartTranscriptionJobInput{
				TranscriptionJobName: jobName,
				LanguageCode:         aws.String(transcribeservice.LanguageCodeEnUs),
				Media: &transcribeservice.Media{
					MediaFileUri: aws.String(config.Source.S3Uri()),
				},
				OutputBucketName: aws.String(config.OutputBucket),
				OutputKey:        aws.String(strings.Replace(config.Source.Key, "."+format, ".json", 1)),
				Settings:         settings,
			})
		if err != nil {
			return errors.Wrap(err, "StartStreamTranscriptionWithContext")
		}
		start = resp.TranscriptionJob.CreationTime
		jobLog.Info("started transcription job", zap.String("format", format))
	} else if err != nil {
		return errors.Wrap(err, "GetTranscriptionJobWithContext")
	} else if existing.TranscriptionJob.CompletionTime != nil {
		log.Info("transcription is complete",
			zap.String("uri", *existing.TranscriptionJob.Transcript.TranscriptFileUri))
		return nil
	} else {
		start = existing.TranscriptionJob.CreationTime
		jobLog.Info(
			"transcription job is already running",
			base.Elapsed(*start))
	}
	if !follow {
		return nil
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
				log.Info("transcription is complete",
					zap.String("uri", *existing.TranscriptionJob.Transcript.TranscriptFileUri))
				return nil
			}
			jobLog.Info(
				"transcription job is still in progress",
				base.Elapsed(*start))
		}
	}
}
