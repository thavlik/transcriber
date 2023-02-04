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
	inputKey := key(input)
	result, err := m.redis.Get(
		ctx,
		inputKey,
	).Result()
	if err == redis.Nil {
		if m.underlying != nil {
			// try and get the value from the underlying storage
			isDisease, err := m.underlying.IsDisease(ctx, input)
			if err != nil {
				return false, err
			}
			// cache the value in redis
			if _, err := m.redis.Set(
				ctx,
				inputKey,
				encode(isDisease),
				0,
			).Result(); err != nil {
				return false, errors.Wrap(err, "redis")
			}
			return isDisease, nil
		}
		return false, diseasecache.ErrNotFound
	} else if err != nil {
		return false, errors.Wrap(err, "redis")
	}
	switch result {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, errors.Errorf("invalid value '%s'", result)
	}
}
