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
	inputBucket  string
	outputBucket string
}

var batchTranscribeCmd = &cobra.Command{
	Use:  "batch-transcribe",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sess := base.AWSSession()
		s3Client := s3.New(sess)
		var nextToken *string
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
				if err == nil {
					if head.LastModified.After(*obj.LastModified) {
						// output file is newer than input file, skip
						base.DefaultLog.Info(
							"transcript already exists, skipping",
							zap.String("key", outKey))
						continue
					}
				}
				if err := batch.TranscribeBatchDefault(
					cmd.Context(),
					&batch.BatchTranscribeSource{
						Bucket: batchTranscribeArgs.inputBucket,
						Key:    key,
					},
					batchTranscribeArgs.outputBucket,
					base.DefaultLog,
				); err != nil {
					return errors.Wrap(err, "transcribe.TranscribeBatchDefault")
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
		"judyscripts",
		"output s3 bucket",
	)
}
