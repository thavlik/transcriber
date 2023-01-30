package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func Search(
	ctx context.Context,
	input string,
	endpoint string,
	subscriptionKey string,
	count int,
	offset int,
) ([]*Image, error) {
	// this was not needed:
	// https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/cognitiveservices/armcognitiveservices
	endpoint += "v7.0/images/search"
	endpoint += fmt.Sprintf("?q=%s", url.QueryEscape(input))
	endpoint += fmt.Sprintf("&count=%d", count)
	if offset != 0 {
		endpoint += fmt.Sprintf("&offset=%d", offset)
	}
	req, err := http.NewRequest(
		"GET",
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // attempt to read error
		return nil, fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, string(body))
	}
	result := &searchResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}
	return convert(result), nil
}
