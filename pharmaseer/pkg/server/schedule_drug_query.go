package server

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) scheduleDrugQuery(
	query string,
	force bool,
) error {
	body, err := json.Marshal(&entity{
		Type:  drugEntity,
		Query: query,
		Force: force,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added drug query to scheduler", zap.String("query", query))
	return nil
}
