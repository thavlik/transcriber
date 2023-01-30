package comprehend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehendmedical"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/transcriber/pkg/util"
	"go.uber.org/zap"
)

type Entity struct {
	Text  string  `json:"text"`
	Type  string  `json:"type"`
	Score float64 `json:"score"`
}

func Comprehend(
	ctx context.Context,
	text string,
	filter []string,
	log *zap.Logger,
) ([]*Entity, error) {
	svc := comprehendmedical.New(util.AWSSession())
	resp, err := svc.DetectEntitiesV2WithContext(
		ctx,
		&comprehendmedical.DetectEntitiesV2Input{
			Text: aws.String(text),
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect entities")
	}
	var entities []*Entity
	for _, entity := range resp.Entities {
		ty := aws.StringValue(entity.Type)
		if filter != nil && !contains(filter, ty) {
			continue
		}
		entities = append(entities, &Entity{
			Text:  aws.StringValue(entity.Text),
			Type:  ty,
			Score: aws.Float64Value(entity.Score),
		})
	}
	return entities, nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
