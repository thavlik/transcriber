package redis_disease_cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
)

type redisDiseaseCache struct {
	redis      *redis.Client
	underlying diseasecache.DiseaseCache
}

func NewRedisDiseaseCache(
	redis *redis.Client,
	underlying diseasecache.DiseaseCache,
) diseasecache.DiseaseCache {
	return &redisDiseaseCache{
		redis:      redis,
		underlying: underlying,
	}
}

func key(input string) string {
	return "disease:" + input
}
