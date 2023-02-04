package server

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) scheduleDrugQuery(
	id string,
) error {
	s.log.Debug("asynchronously querying drug details", zap.String("id", id))
	body, err := json.Marshal(&entity{
		Type: drug,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added drug query to scheduler", zap.String("id", id))
	return nil
}
