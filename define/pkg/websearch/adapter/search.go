package adapter

import (
	"context"

	"github.com/thavlik/transcriber/define/pkg/websearch"
	bing_websearch "github.com/thavlik/transcriber/define/pkg/websearch/bing"

	"github.com/pkg/errors"
)

type Service string

const (
	Bing Service = "bing"
)

func Search(
	ctx context.Context,
	service Service,
	query string,
	endpoint string,
	subscriptionKey string,
	count int,
	offset int,
) ([]*websearch.Result, error) {
	switch service {
	case Bing:
		return bing_websearch.Search(
			ctx,
			query,
			endpoint,
			subscriptionKey,
			count,
			offset,
		)
	default:
		return nil, errors.Errorf("unsupported service '%s'", service)
	}
}
