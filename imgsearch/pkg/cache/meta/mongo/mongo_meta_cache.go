package mongo

import (
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/meta"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoMetaCache struct {
	c *mongo.Collection
}

func NewMongoMetaCache(c *mongo.Collection) meta.ImageMetaCache {
	return &mongoMetaCache{c}
}
