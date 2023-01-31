package cache

import (
	"context"
	"io"
)

// Get retrieves the image data from the cache.
func (c *ImageCache) Get(
	ctx context.Context,
	hash string,
) (io.ReadCloser, error) {
	return c.dataCache.Get(
		ctx,
		hash,
	)
}
