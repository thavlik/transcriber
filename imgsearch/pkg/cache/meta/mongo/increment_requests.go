package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (c *mongoMetaCache) IncrementRequests(
	ctx context.Context,
	hash string,
) error {
	if _, err := c.c.UpdateOne(
		ctx,
		map[string]interface{}{
			"_id": hash,
		},
		map[string]interface{}{
			"$inc": map[string]interface{}{
				"requests": 1,
			},
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
