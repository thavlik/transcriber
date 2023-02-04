package redis

import (
	"context"

	"github.com/pkg/errors"
)

func (s *redisScheduler) Add(entities ...string) error {
	p := s.redis.Pipeline()
	for _, entity := range entities {
		p.ZIncrBy(
			context.Background(),
			s.key,
			1.0,
			entity,
		)
	}
	p.Publish(
		context.Background(),
		channelName(s.key),
		"1",
	)
	if _, err := p.Exec(context.Background()); err != nil {
		return errors.Wrap(err, "redis.Exec")
	}
	return nil
}
