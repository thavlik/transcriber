package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
)

type s3DataCache struct {
	sess   *session.Session
	s3     *s3.S3
	bucket string
}

func NewS3DataCache(bucket string) data.ImageDataCache {
	sess := base.AWSSession()
	return &s3DataCache{
		sess,
		s3.New(sess),
		bucket,
	}
}
