package memory

import (
	"context"

	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

func (m *memoryEntityCache) Cache(
	ctx context.Context,
	entity *comprehend.Entity,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.l <- struct{}{}:
		defer func() { <-m.l }()
		m.m[entity.Hash()] = entity
		return nil
	}
}
