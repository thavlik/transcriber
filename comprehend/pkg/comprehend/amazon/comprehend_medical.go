package amazon_comprehend

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehendmedical"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

// ComprehendMedical detects entities in medical text.
// This is backed by Amazon's Comprehend Medical service.
// If includeTypes is not empty, the detected entity type must be in the list, otherwise it is filtered.
// If excludeTypes is not empty, the detected entity type must not be in the list, otherwise it is filtered.
func ComprehendMedical(
	ctx context.Context,
	text string,
	filter *comprehend.Filter,
) ([]*comprehend.Entity, error) {
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
		filter,
	), nil
}

func convertMedicalEntities(
	entities []*comprehendmedical.Entity,
	filter *comprehend.Filter,
) (result []*comprehend.Entity) {
	for _, entity := range entities {
		e := &comprehend.Entity{
			Text:  strings.ToLower(aws.StringValue(entity.Text)),
			Type:  aws.StringValue(entity.Type),
			Score: aws.Float64Value(entity.Score),
		}
		if filter.Matches(e) {
			result = append(result, e)
		}
	}
	return
}
