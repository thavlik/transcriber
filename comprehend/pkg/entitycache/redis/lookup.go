package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
)

func (m *redisEntityCache) Lookup(
	ctx context.Context,
	hash string,
) (*comprehend.Entity, error) {
	body, err := m.redis.Get(
		ctx,
		key(hash),
	).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, entitycache.ErrNotFound
		}
		return nil, errors.Wrap(err, "redis")
	}
	entity := new(comprehend.Entity)
	if err := json.Unmarshal([]byte(body), entity); err != nil {
		return nil, errors.Wrap(err, "json")
	}
	return entity, nil
}
