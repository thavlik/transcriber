package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
)

func (m *mongoHistory) Push(
	ctx context.Context,
	search *history.Search,
) error {
	_, err := m.c.InsertOne(
		ctx,
		search.AsMongoMap(),
	)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
