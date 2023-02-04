package s3_pdbcache

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	"go.uber.org/zap"
)

type s3PDBCache struct {
	sess       *session.Session
	s3         *s3.S3
	bucketName string
	log        *zap.Logger
}

func NewS3PDBCache(
	bucketName string,
	log *zap.Logger,
) pdbcache.PDBCache {
	sess := base.AWSSession()
	return &s3PDBCache{
		sess,
		s3.New(sess),
		bucketName,
		log,
	}
}

func pdbKey(id string) string {
	return fmt.Sprintf("%s.pdb", id)
}
