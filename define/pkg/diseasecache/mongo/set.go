package mongo_disease_cache

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mongoDiseaseCache) Set(
	ctx context.Context,
	input string,
	isDisease bool,
) error {
	if _, err := m.c.UpdateOne(
		ctx,
		map[string]interface{}{
			"_id": input,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"is_disease": isDisease,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
