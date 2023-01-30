package search

func convert(result *searchResult) []*Image {
	images := make([]*Image, len(result.Value))
	for i, v := range result.Value {
		images[i] = &Image{
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
