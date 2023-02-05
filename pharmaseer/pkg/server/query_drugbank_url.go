package server

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/websearch/adapter"
)

var errFailedToFindDrugBankURL = errors.New("failed to find drugbank url")

func queryDrugBankURL(
	ctx context.Context,
	query string,
	service adapter.Service,
	endpoint string,
	subscriptionKey string,
) (string, error) {
	results, err := adapter.Search(
		ctx,
		service,
		"go.drugbank.com "+query,
		endpoint,
		subscriptionKey,
		10,
		0,
	)
	if err != nil {
		return "", errors.Wrap(err, "search failed")
	}
	for _, result := range results {
		if strings.HasPrefix(result.URL, "https://go.drugbank.com/drugs/") {
			return result.URL, nil
		}
	}
	return "", errFailedToFindDrugBankURL
}
