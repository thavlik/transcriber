package redis

import (
	"time"

	"github.com/bsm/redislock"
	redis "github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
)

type redisScheduler struct {
	redis   *redis.Client
	locker  *redislock.Client
	key     string
	lockTTL time.Duration
}

func NewRedisScheduler(
	redisClient *redis.Client,
	key string,
	lockTTL time.Duration,
) scheduler.Scheduler {
	return &redisScheduler{
		redisClient,
		redislock.New(redisClient),
		key,
		lockTTL,
	}
}

// channelName redis publish/subscribe channel name
func channelName(key string) string {
	return key + ":ch"
}
