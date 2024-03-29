package bing_websearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/websearch"
)

func Search(
	ctx context.Context,
	query string,
	endpoint string,
	subscriptionKey string,
	count int,
	offset int,
) ([]*websearch.Result, error) {
	endpoint += "v7.0/search"
	endpoint += fmt.Sprintf("?q=%s", url.QueryEscape(query))
	endpoint += fmt.Sprintf("&count=%d", count)
	if offset != 0 {
		endpoint += fmt.Sprintf("&offset=%d", offset)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)
	resp, err := (&http.Client{
		Timeout: 20 * time.Second,
	}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // attempt to read error
		return nil, fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, string(body))
	}
	result := &webSearchResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}
	return convert(result), nil
}
