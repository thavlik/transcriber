package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func (s *s3ThumbCache) Set(
	id string,
	r io.Reader,
) (err error) {
	key := thumbKey(id)
	log := s.log.With(
		zap.String("bucket", s.bucketName),
		zap.String("key", key))
	if _, err := s3manager.NewUploader(base.AWSSession()).
		Upload(&s3manager.UploadInput{
			Bucket: aws.String(s.bucketName),  // all image uploads go into one bucket
			Key:    aws.String(key),           // videoID.webm
			Body:   r,                         // videoID.webm
			ACL:    aws.String("public-read"), // "x-amz-acl" https://docs.digitalocean.com/reference/api/spaces-api/
		}); err != nil {
		return errors.Wrap(err, "s3")
	}
	log.Debug("completed multipart upload")
	return nil
}
