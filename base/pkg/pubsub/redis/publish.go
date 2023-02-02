package redis_pubsub

import (
	"context"

	"github.com/pkg/errors"
)

func (p *redisPubSub) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
) error {
	if _, err := p.redis.Publish(
		ctx,
		topic,
		payload,
	).Result(); err != nil {
		return errors.Wrap(err, "redis")
	}
	return nil
}
