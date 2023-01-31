package cache

import (
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/meta"
)

type ImageCache struct {
	metaCache meta.ImageMetaCache
	dataCache data.ImageDataCache
}

func NewImageCache(
	metaCache meta.ImageMetaCache,
	dataCache data.ImageDataCache,
) *ImageCache {
	return &ImageCache{
		metaCache: metaCache,
		dataCache: dataCache,
	}
}
