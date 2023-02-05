package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

func (s *s3ThumbCache) Has(
	ctx context.Context,
	id string,
) (bool, error) {
	head, err := s.s3.HeadObjectWithContext(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(thumbKey(id)),
		})
	if err != nil {
		// HeadObject returns non-standard 404 https://github.com/aws/aws-sdk-go/issues/2095
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, errors.Wrap(err, "s3")
	}
	return aws.Int64Value(head.ContentLength) > 0, nil
}
