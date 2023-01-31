package mongo

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

func (c *mongoMetaCache) Set(
	ctx context.Context,
	img *search.Image,
) error {
	// simple hack for converting a struct to a map and
	// then assigning the hash as the mongo _id field
	body, err := json.Marshal(img)
	if err != nil {
		return errors.Wrap(err, "json")
	}
	doc := make(map[string]interface{})
	if err := json.Unmarshal(body, &doc); err != nil {
		return errors.Wrap(err, "json")
	}
	doc["_id"] = img.Hash()
	if _, err := c.c.InsertOne(
		ctx,
		doc,
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
