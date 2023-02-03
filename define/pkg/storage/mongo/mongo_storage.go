package mongo

import (
	"github.com/thavlik/transcriber/define/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStorage struct {
	c *mongo.Collection
}

func NewMongoStorage(c *mongo.Collection) storage.Storage {
	return &mongoStorage{c}
}
