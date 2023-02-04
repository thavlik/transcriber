package amazon_comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	awscomprehend "github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

// Comprehend detects entities in text.
// This is the standard Comprehend service.
// See https://docs.aws.amazon.com/comprehend/latest/dg/how-entities.html
// Use the filter parameter for fine-grain control over which entities are returned.
func Comprehend(
	ctx context.Context,
	text string,
	languageCode string,
	filter *comprehend.Filter,
) ([]*comprehend.Entity, error) {
	svc := awscomprehend.New(base.AWSSession())
	resp, err := svc.DetectEntitiesWithContext(
		ctx,
		&awscomprehend.DetectEntitiesInput{
			Text:         aws.String(text),
			LanguageCode: aws.String(languageCode),
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
	entities []*awscomprehend.Entity,
	filter *comprehend.Filter,
) (result []*comprehend.Entity) {
	for _, entity := range entities {
		e := &comprehend.Entity{
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
