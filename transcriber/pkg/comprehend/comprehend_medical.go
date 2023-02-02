package comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehendmedical"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

// ComprehendMedical detects entities in medical text.
// This is backed by Amazon's Comprehend Medical service.
// If includeTypes is not empty, the detected entity type must be in the list, otherwise it is filtered.
// If excludeTypes is not empty, the detected entity type must not be in the list, otherwise it is filtered.
func ComprehendMedical(
	ctx context.Context,
	text string,
	includeTypes []string,
	excludeTypes []string,
	log *zap.Logger,
) ([]*Entity, error) {
	svc := comprehendmedical.New(base.AWSSession())
	resp, err := svc.DetectEntitiesV2WithContext(
		ctx,
		&comprehendmedical.DetectEntitiesV2Input{
			Text: aws.String(text),
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect entities")
	}
	return convertMedicalEntities(
		resp.Entities,
		includeTypes,
		excludeTypes,
	), nil
}

func convertMedicalEntities(
	entities []*comprehendmedical.Entity,
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
