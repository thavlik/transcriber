package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (s *redisScheduler) List() ([]string, error) {
	entities, err := s.redis.ZRevRange(
		context.Background(),
		s.key,
		0,
		-1,
	).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis.ZRevRange")
	}
	return entities, nil
}
