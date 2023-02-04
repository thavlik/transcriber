package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"
)

func (c *mongoMetaCache) Set(
	ctx context.Context,
	img *imgsearch.Image,
	fileHash string,
) error {
	doc, err := img.AsMap()
	if err != nil {
		return err
	}
	doc["_id"] = img.Hash()
	doc["fileHash"] = fileHash
	if _, err := c.c.InsertOne(
		ctx,
		doc,
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
