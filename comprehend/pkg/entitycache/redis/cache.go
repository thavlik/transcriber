package redis

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

func (m *redisEntityCache) Cache(
	ctx context.Context,
	entity *comprehend.Entity,
) error {
	body, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	if _, err := m.redis.Set(
		ctx,
		key(entity.Hash()),
		string(body),
		0,
	).Result(); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
