package mongo_disease_cache

import (
	"context"

	"github.com/pkg/errors"
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
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "mongo")
	}
	return doc["is_disease"] == true, nil
}
