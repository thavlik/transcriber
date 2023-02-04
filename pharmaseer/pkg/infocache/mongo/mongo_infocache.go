package mongo

import (
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoInfoCache struct {
	c *mongo.Collection
}

func NewMongoInfoCache(db *mongo.Database) infocache.InfoCache {
	return &mongoInfoCache{
		c: db.Collection("drugs"),
	}
}
