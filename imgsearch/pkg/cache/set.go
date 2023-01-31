package cache

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

// Set stores the image data in the cache.
func (c *ImageCache) Set(
	ctx context.Context,
	img *search.Image,
	r io.Reader,
) error {
	// cache the image bytes in the storage backend
	if err := c.dataCache.Set(ctx, img, r); err != nil {
		return errors.Wrap(err, "datacache")
	}
	// cache the image metadata in the database
	if err := c.metaCache.Set(ctx, img); err != nil {
		return errors.Wrap(err, "metacache")
	}
	return nil
}
