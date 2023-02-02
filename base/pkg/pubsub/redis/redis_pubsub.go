package redis_pubsub

import (
	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"go.uber.org/zap"
)

type redisPubSub struct {
	redis *redis.Client
	log   *zap.Logger
}

func NewRedisPubSub(
	redis *redis.Client,
	log *zap.Logger,
) pubsub.PubSub {
	return &redisPubSub{
		redis,
		log,
	}
}
