package comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

// Comprehend detects entities in text.
// This is the standard Comprehend service.
// See https://docs.aws.amazon.com/comprehend/latest/dg/how-entities.html
// Use the filter parameter for fine-grain control over which entities are returned.
func Comprehend(
	ctx context.Context,
	text string,
	filter *Filter,
	log *zap.Logger,
) ([]*Entity, error) {
	svc := comprehend.New(base.AWSSession())
	resp, err := svc.DetectEntitiesWithContext(
		ctx,
		&comprehend.DetectEntitiesInput{
			Text:         aws.String(text),
			LanguageCode: aws.String("en"),
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect entities")
	}
	return convertEntities(
		resp.Entities,
		filter,
	), nil
}

func convertEntities(
	entities []*comprehend.Entity,
	filter *Filter,
) (result []*Entity) {
	for _, entity := range entities {
		e := &Entity{
			Text:  aws.StringValue(entity.Text),
			Type:  aws.StringValue(entity.Type),
			Score: aws.Float64Value(entity.Score),
		}
		if filter.Matches(e) {
			result = append(result, e)
		}
	}
	return
}
