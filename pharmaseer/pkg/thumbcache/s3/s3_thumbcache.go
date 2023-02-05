package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/pharmaseer/pkg/thumbcache"
	"go.uber.org/zap"
)

type s3ThumbCache struct {
	sess       *session.Session
	s3         *s3.S3
	bucketName string
	log        *zap.Logger
}

func NewS3ThumbCache(bucketName string, log *zap.Logger) thumbcache.ThumbCache {
	sess := base.AWSSession()
	return &s3ThumbCache{
		sess,
		s3.New(sess),
		bucketName,
		log,
	}
}

func thumbKey(id string) string {
	return fmt.Sprintf("%s.svg", id)
}
