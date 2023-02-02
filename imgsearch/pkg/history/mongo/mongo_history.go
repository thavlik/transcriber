package mongo

import (
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoHistory struct {
	c *mongo.Collection
}

func NewMongoHistory(c *mongo.Collection) history.History {
	return &mongoHistory{
		c: c,
	}
}
