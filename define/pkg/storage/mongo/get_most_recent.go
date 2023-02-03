package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mongoStorage) GetMostRecent(
	ctx context.Context,
	input string,
) (*storage.Definition, error) {
	result := m.c.FindOne(
		ctx,
		map[string]interface{}{
			"input": input,
		},
		options.FindOne().SetSort(
			map[string]interface{}{
				"timestamp": -1,
			},
		),
	)
	var raw struct {
		ID        string `bson:"_id"`
		Input     string `bson:"input"`
		Output    string `bson:"output"`
		Timestamp int64  `bson:"timestamp"`
	}
	if err := result.Decode(&raw); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, storage.ErrNotCached
		}
		return nil, errors.Wrap(err, "mongo")
	}
	return &storage.Definition{
		ID:        raw.ID,
		Input:     raw.Input,
		Output:    raw.Output,
		Timestamp: time.Unix(0, raw.Timestamp),
	}, nil
}
