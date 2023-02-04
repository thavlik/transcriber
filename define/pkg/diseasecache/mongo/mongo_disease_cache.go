package mongo_disease_cache

import (
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDiseaseCache struct {
	c *mongo.Collection
}

func NewMongoDiseaseCache(
	c *mongo.Collection,
) diseasecache.DiseaseCache {
	return &mongoDiseaseCache{c}
}
