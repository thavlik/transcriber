package redis_disease_cache

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
)

func (m *redisDiseaseCache) IsDisease(
	ctx context.Context,
	input string,
) (bool, error) {
	result, err := m.redis.Get(
		ctx,
		key(input),
	).Result()
	if err == redis.Nil {
		if m.underlying != nil {
			// try and get the value from the underlying storage
			return m.underlying.IsDisease(ctx, input)
		}
		return false, diseasecache.ErrNotFound
	} else if err != nil {
		return false, errors.Wrap(err, "redis")
	}
	return result == "1", nil
}
