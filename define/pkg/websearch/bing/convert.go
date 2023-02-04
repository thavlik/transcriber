package bing_websearch

import "github.com/thavlik/transcriber/define/pkg/websearch"

func convert(result *webSearchResult) []*websearch.Result {
	results := make([]*websearch.Result, len(result.WebPages.Value))
	for i, v := range result.WebPages.Value {
		results[i] = &websearch.Result{
			Name:             v.Name,
			URL:              v.URL,
			IsFamilyFriendly: v.IsFamilyFriendly,
		}
	}
	return results
}
