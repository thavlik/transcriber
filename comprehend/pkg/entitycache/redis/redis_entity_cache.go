package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
)

type redisEntityCache struct {
	redis *redis.Client
}

func NewRedisEntityCache(
	redis *redis.Client,
) entitycache.EntityCache {
	return &redisEntityCache{redis}
}

func key(hash string) string {
	return "ent:" + hash
}
