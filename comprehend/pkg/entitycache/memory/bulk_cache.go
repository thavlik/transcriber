package memory

import (
	"context"

	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

func (m *memoryEntityCache) BulkCache(
	ctx context.Context,
	entities []*comprehend.Entity,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.l <- struct{}{}:
		defer func() { <-m.l }()
		for _, entity := range entities {
			m.m[entity.Hash()] = entity
		}
		return nil
	}
}
