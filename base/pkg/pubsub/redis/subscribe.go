package redis_pubsub

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/pubsub"
)

func (p *redisPubSub) Subscribe(
	ctx context.Context,
	topic string,
) (pubsub.Subscription, error) {
	return &redisSubscription{
		redis: p.redis,
		stop:  make(chan struct{}, 1),
		topic: topic,
		log:   p.log,
	}, nil
}
