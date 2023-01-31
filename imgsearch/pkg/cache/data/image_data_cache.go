package data

import (
	"context"
	"errors"
	"io"

	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

var ErrNotCached = errors.New("not cached")

type ImageDataCache interface {
	Get(ctx context.Context, hash string) (io.ReadCloser, error)
	Set(context.Context, *search.Image, io.Reader) error
}
