package bing_search

import "github.com/thavlik/transcriber/imgsearch/pkg/search"

func convert(result *searchResult) []*search.Image {
	images := make([]*search.Image, len(result.Value))
	for i, v := range result.Value {
		images[i] = &search.Image{
			ContentURL:     v.ContentUrl,
			ContentSize:    v.ContentSize,
			ThumbnailURL:   v.ThumbnailUrl,
			EncodingFormat: v.EncodingFormat,
			HostPageURL:    v.HostPageUrl,
			Width:          v.Width,
			Height:         v.Height,
			AccentColor:    v.AccentColor,
		}
	}
	return images
}
