package meta

import (
	"context"

	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

// ImageMetaCache is a cache for image metadata.
// It is used to keep track of which images have been cached.
// This is never read by the application, only written to.
// It is utilized so the source for an image can be queried.
type ImageMetaCache interface {
	Set(context.Context, *search.Image) error
	Increment(ctx context.Context, hash string) error
}
