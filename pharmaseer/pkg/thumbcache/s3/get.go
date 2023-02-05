package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
)

func (s *s3ThumbCache) Get(
	ctx context.Context,
	videoID string,
	w io.Writer,
) error {
	result, err := s.s3.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(thumbKey(videoID)),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return api.ErrNotCached
		}
		return errors.Wrap(err, "s3")
	}
	if _, err := io.Copy(w, result.Body); err != nil {
		return errors.Wrap(err, "copy")
	}
	return nil
}
