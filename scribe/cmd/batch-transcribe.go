package main

import (
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/scribe/pkg/batch"
	"go.uber.org/zap"
)

var batchTranscribeArgs struct {
	inputBucket      string
	outputBucket     string
	maxSpeakerLabels int64
	limit            int64
	follow           bool
}

var batchTranscribeCmd = &cobra.Command{
	Use:  "batch-transcribe",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sess := base.AWSSession()
		s3Client := s3.New(sess)
		var nextToken *string
		var transcribedCount int64
		for {
			resp, err := s3Client.ListObjectsWithContext(
				cmd.Context(),
				&s3.ListObjectsInput{
					Bucket: &batchTranscribeArgs.inputBucket,
					Marker: nextToken,
				},
			)
			if err != nil {
				return errors.Wrap(err, "s3Client.ListObjectsWithContext")
			}
			for _, obj := range resp.Contents {
				key := aws.StringValue(obj.Key)
				ext := filepath.Ext(key)
				if ext != ".mp3" {
					continue
				}
				outKey := key[:len(key)-len(ext)] + ".json"
				head, err := s3Client.HeadObjectWithContext(
					cmd.Context(),
					&s3.HeadObjectInput{
						Bucket: &batchTranscribeArgs.outputBucket,
						Key:    aws.String(outKey),
					},
				)
				if err == nil &&
					aws.Int64Value(head.ContentLength) > 0 &&
					head.LastModified.After(*obj.LastModified) {
					base.DefaultLog.Info(
						"transcript already exists, skipping",
						zap.String("bucket", batchTranscribeArgs.outputBucket),
						zap.String("key", outKey))
					continue
				}
				if err := batch.TranscribeBatchDefault(
					cmd.Context(),
					&batch.Config{
						Source: &batch.BatchTranscribeSource{
							Bucket: batchTranscribeArgs.inputBucket,
							Key:    key,
						},
						OutputBucket:     batchTranscribeArgs.outputBucket,
						MaxSpeakerLabels: batchTranscribeArgs.maxSpeakerLabels,
					},
					batchTranscribeArgs.follow,
					base.DefaultLog,
				); err != nil {
					return errors.Wrap(err, "transcribe.TranscribeBatchDefault")
				}
				transcribedCount++
				if batchTranscribeArgs.limit > 0 && transcribedCount >= batchTranscribeArgs.limit {
					base.DefaultLog.Info(
						"limit reached, exiting",
						zap.Int64("limit", batchTranscribeArgs.limit))
					return nil
				}
			}
			if resp.NextMarker == nil {
				break
			}
			nextToken = resp.NextMarker
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchTranscribeCmd)
	batchTranscribeCmd.Flags().StringVarP(
		&batchTranscribeArgs.inputBucket,
		"input-bucket",
		"i",
		"judymp3",
		"input s3 bucket",
	)
	batchTranscribeCmd.Flags().StringVarP(
		&batchTranscribeArgs.outputBucket,
		"output-bucket",
		"o",
		"judyscripts-v2",
		"output s3 bucket",
	)
	batchTranscribeCmd.Flags().Int64VarP(
		&batchTranscribeArgs.maxSpeakerLabels,
		"max-speaker-labels",
		"m",
		10,
		"maximum number of speaker labels",
	)
	batchTranscribeCmd.Flags().Int64VarP(
		&batchTranscribeArgs.limit,
		"limit",
		"l",
		0,
		"limit number of files to transcribe (0 = no limit)",
	)
	batchTranscribeCmd.Flags().BoolVarP(
		&batchTranscribeArgs.follow,
		"follow",
		"f",
		false,
		"follow each transcription job until it completes",
	)
}
