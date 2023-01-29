package comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pkg/util"
	"go.uber.org/zap"
)

type Entity struct {
	Text  string
	Type  string
	Score float64
}

func Comprehend(
	ctx context.Context,
	text string,
	log *zap.Logger,
) ([]*Entity, error) {
	svc := comprehend.New(util.AWSSession())
	resp, err := svc.DetectEntitiesWithContext(
		ctx,
		&comprehend.DetectEntitiesInput{
			LanguageCode: aws.String("en"),
			Text:         aws.String(text),
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect entities")
	}
	entities := make([]*Entity, len(resp.Entities))
	for i, entity := range resp.Entities {
		entities[i] = &Entity{
			Text:  aws.StringValue(entity.Text),
			Type:  aws.StringValue(entity.Type),
			Score: aws.Float64Value(entity.Score),
		}
	}
	return entities, nil
}
