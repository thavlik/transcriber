package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (s *redisScheduler) Remove(entity string) error {
	if _, err := s.redis.ZRem(
		context.Background(),
		s.key,
		entity,
	).Result(); err != nil {
		return errors.Wrap(err, "redis.ZRem")
	}
	return nil
}
