package mongo_disease_cache

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *mongoDiseaseCache) IsDisease(
	ctx context.Context,
	input string,
) (bool, error) {
	doc := make(map[string]interface{})
	if err := m.c.FindOne(
		ctx,
		map[string]interface{}{
			"_id": input,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return false, diseasecache.ErrNotFound
	} else if err != nil {
		return false, errors.Wrap(err, "mongo")
	}
	return doc["is_disease"] == true, nil
}
