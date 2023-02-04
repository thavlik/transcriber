package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"
)

// Set stores the image data in the cache.
func (c *ImageCache) Set(
	ctx context.Context,
	img *imgsearch.Image,
	r io.Reader,
) error {
	// calculate the hash of the file as we read it
	h := md5.New()
	r = io.TeeReader(r, h)

	// cache the image bytes in the storage backend
	if err := c.dataCache.Set(
		ctx,
		img,
		r,
	); err != nil {
		return errors.Wrap(err, "datacache")
	}

	// cache the image metadata in the database
	raw := h.Sum(nil)
	fileHash := hex.EncodeToString(raw[:])
	if err := c.metaCache.Set(
		ctx,
		img,
		fileHash,
	); err != nil {
		return errors.Wrap(err, "metacache")
	}

	return nil
}
