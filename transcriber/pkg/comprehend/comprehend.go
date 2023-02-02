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
// If includeTypes is not empty, the detected entity type must be in the list, otherwise it is filtered.
// If excludeTypes is not empty, the detected entity type must not be in the list, otherwise it is filtered.
func Comprehend(
	ctx context.Context,
	text string,
	includeTypes []string,
	excludeTypes []string,
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
		includeTypes,
		excludeTypes,
	), nil
}

func convertEntities(
	entities []*comprehend.Entity,
	includeTypes []string,
	excludeTypes []string,
) (result []*Entity) {
	for _, entity := range entities {
		ty := aws.StringValue(entity.Type)
		if filter(ty, includeTypes, excludeTypes) {
			continue
		}
		result = append(result, &Entity{
			Text:  aws.StringValue(entity.Text),
			Type:  ty,
			Score: aws.Float64Value(entity.Score),
		})
	}
	return
}
