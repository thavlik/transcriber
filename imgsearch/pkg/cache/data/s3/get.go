package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"

	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
)

func (c *s3DataCache) Get(
	ctx context.Context,
	hash string,
) (io.ReadCloser, error) {
	result, err := c.s3.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(hash),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return nil, data.ErrNotCached
		}
		return nil, errors.Wrap(err, "s3")
	}
	return result.Body, nil
}
