package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/storage"
)

func (m *mongoStorage) Insert(
	ctx context.Context,
	def *storage.Definition,
) error {
	if _, err := m.c.InsertOne(
		ctx,
		map[string]interface{}{
			"_id":     def.ID,
			"input":   def.Input,
			"output":  def.Output,
			"updated": def.Timestamp.UnixNano(),
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
