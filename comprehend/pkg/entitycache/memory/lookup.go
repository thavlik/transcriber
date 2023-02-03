package memory

import (
	"context"

	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
)

func (m *memoryEntityCache) Lookup(
	ctx context.Context,
	hash string,
) (*comprehend.Entity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m.l <- struct{}{}:
		defer func() { <-m.l }()
		if e, ok := m.m[hash]; ok {
			return e, nil
		}
		return nil, entitycache.ErrNotFound
	}
}
