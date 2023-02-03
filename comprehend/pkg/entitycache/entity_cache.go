package entitycache

import (
	"context"
	"errors"

	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

var ErrNotFound = errors.New("entity not found")

type EntityCache interface {
	Cache(context.Context, *comprehend.Entity) error
	BulkCache(context.Context, []*comprehend.Entity) error
	Lookup(context.Context, string) (*comprehend.Entity, error)
}
