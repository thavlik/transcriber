package bing_imgsearch

import "github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"

func convert(result *searchResult) []*imgsearch.Image {
	images := make([]*imgsearch.Image, len(result.Value))
	for i, v := range result.Value {
		images[i] = &imgsearch.Image{
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
