package thumbcache

import (
	"context"
	"io"
)

type ThumbCache interface {
	Has(ctx context.Context, id string) (bool, error)
	Del(id string) error
	Set(id string, r io.Reader) error
	Get(ctx context.Context, id string, w io.Writer) error
}
