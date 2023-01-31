package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"

	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

func (c *s3DataCache) Set(
	ctx context.Context,
	img *search.Image,
	r io.Reader,
) error {
	if _, err := s3manager.NewUploader(c.sess).
		UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(EncodeHash(img.Hash())),
			Body:   r,
			ACL:    aws.String("public-read"),
		}); err != nil {
		// Warning: there may be a dangling multipart upload.
		// This is bad because s3 charges for the storage they consume,
		// even if the object is never visible. We can clean them up,
		// but in the event multiple users are uploading the same file
		// at the same time, we may end up deleting a valid upload.
		// This is a rare case, so we'll just ignore it for now.
		// The workaround is manually free multipart uploads on a
		// regular basis.
		return errors.Wrap(err, "s3")
	}
	return nil
}
