package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (c *mongoInfoCache) HasDrug(
	ctx context.Context,
	query string,
) (bool, error) {
	n, err := c.c.CountDocuments(
		ctx,
		map[string]interface{}{
			"_id": query,
		},
	)
	if err != nil {
		return false, errors.Wrap(err, "mongo")
	}
	return n > 0, nil
}
