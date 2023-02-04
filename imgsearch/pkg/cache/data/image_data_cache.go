package data

import (
	"context"
	"errors"
	"io"

	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"
)

var ErrNotCached = errors.New("not cached")

type ImageDataCache interface {
	Get(ctx context.Context, hash string) (io.ReadCloser, error)
	Set(context.Context, *imgsearch.Image, io.Reader) error
}
