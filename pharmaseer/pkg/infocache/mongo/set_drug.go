package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *mongoInfoCache) SetDrug(
	ctx context.Context,
	query string,
	drug *api.DrugDetails,
) error {
	if _, err := c.c.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": query,
		},
		map[string]interface{}{
			"$set": drug.AsMap(),
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
