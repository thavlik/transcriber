package adapter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	amazon_comprehend "github.com/thavlik/transcriber/comprehend/pkg/comprehend/amazon"
)

type Model string

const (
	AmazonComprehend        Model = "amazon-comprehend"
	AmazonComprehendMedical Model = "amazon-comprehend-medical"
)

func Comprehend(
	ctx context.Context,
	model Model,
	text string,
	languageCode string,
	filter *comprehend.Filter,
) ([]*comprehend.Entity, error) {
	switch model {
	case "":
		return nil, errors.New("missing model parameter")
	case AmazonComprehend:
		return amazon_comprehend.Comprehend(
			ctx,
			text,
			languageCode,
			filter,
		)
	case AmazonComprehendMedical:
		return amazon_comprehend.ComprehendMedical(
			ctx,
			text,
			filter,
		)
	default:
		return nil, errors.Errorf("unrecognized model '%s'", model)
	}
}
