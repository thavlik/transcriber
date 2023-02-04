package redis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/base/pkg/base"
)

func (s *redisScheduler) Rand() (string, error) {
	entity, err := s.redis.ZRandMember(
		context.Background(),
		s.key,
		1,
	).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "redis.ZRandMember")
	} else if len(entity) != 1 {
		panic(base.Unreachable)
	}
	return entity[0], nil
}
