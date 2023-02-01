package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

func (s *Server) cacheImage(
	ctx context.Context,
	img *search.Image,
	w io.Writer,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		img.ContentURL,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "newrequest")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != img.ContentType() {
		// sanity check to ensure we're not caching the wrong thing
		// if this fails, it's probably an image with a rare mimetype
		// in any case, we probably don't want to show it to the user
		return fmt.Errorf(
			"content type mismatch: expected %s, got %s",
			img.ContentType(),
			resp.Header.Get("Content-Type"),
		)
	}
	if err := s.imageCache.Set(
		ctx,
		img,
		io.TeeReader(resp.Body, w),
	); err != nil {
		return errors.Wrap(err, "imagecache.Set")
	}
	return nil
}
