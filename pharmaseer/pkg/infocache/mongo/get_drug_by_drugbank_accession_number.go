package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *mongoInfoCache) GetDrugByDrugBankAccessionNumber(
	ctx context.Context,
	drugBankAccessionNumber string,
) (*api.DrugDetails, error) {
	doc := make(map[string]interface{})
	if err := c.c.FindOne(
		ctx,
		map[string]interface{}{
			"drugBankAccessionNumber": drugBankAccessionNumber,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return api.ConvertDrugDetails(doc), nil
}
