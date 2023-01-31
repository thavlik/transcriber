package cache

import (
	"context"
)

// Increment increments the counter keeping track of how many
// times an image with the given hash has been requested.
func (c *ImageCache) IncrementRequests(
	ctx context.Context,
	hash string,
) error {
	return c.metaCache.IncrementRequests(ctx, hash)
}
