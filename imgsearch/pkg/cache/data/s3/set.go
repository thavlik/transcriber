package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"

	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

func (c *s3DataCache) Set(
	ctx context.Context,
	img *search.Image,
	r io.Reader,
) error {
	key := aws.String(EncodeHash(img.Hash()))
	if resp, err := s3manager.NewUploader(c.sess).
		UploadWithContext(
			ctx,
			&s3manager.UploadInput{
				Bucket: aws.String(c.bucket),
				Key:    key,
				Body:   r,
				ACL:    aws.String("public-read"),
			},
		); err != nil {
		// Try and clean up the dangling multipart upload
		_, _ = c.s3.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
			Bucket:   aws.String(c.bucket),
			Key:      key,
			UploadId: aws.String(resp.UploadID),
		})
		return errors.Wrap(err, "s3")
	}
	return nil
}
