package adapter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"
	bing_imgsearch "github.com/thavlik/transcriber/imgsearch/pkg/imgsearch/bing"
)

type SearchService string

const (
	Bing SearchService = "bing"
)

func Search(
	ctx context.Context,
	service SearchService,
	input string,
	endpoint string,
	subscriptionKey string,
	count int,
	offset int,
) (*imgsearch.Result, error) {
	switch service {
	case Bing:
		return bing_imgsearch.Search(
			ctx,
			input,
			endpoint,
			subscriptionKey,
			count,
			offset,
		)
	default:
		return nil, errors.Errorf("unrecognized search service '%s'", service)
	}
}
