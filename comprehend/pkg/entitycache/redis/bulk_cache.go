package redis

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

func (m *redisEntityCache) BulkCache(
	ctx context.Context,
	entities []*comprehend.Entity,
) error {
	p := m.redis.Pipeline()
	for _, entity := range entities {
		body, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		p.Set(
			ctx,
			key(entity.Hash()),
			string(body),
			0,
		)
	}
	if _, err := p.Exec(ctx); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
