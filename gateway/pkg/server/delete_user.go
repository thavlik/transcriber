package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"github.com/thavlik/transcriber/base/pkg/iam/api"
	"go.uber.org/zap"
)

func (s *Server) DeleteUser(ctx context.Context, req api.DeleteUser) (_ *api.Void, err error) {
	if req.ID == "" {
		user, err := s.iam.GetUser(context.Background(), req.Username)
		if err == iam.ErrUserNotFound {
			return &api.Void{}, nil
		} else if err != nil {
			return nil, errors.Wrap(err, "iam.GetUser")
		}
		req.ID = user.ID
	}
	if err := s.iam.DeleteUser(req.Username); err != nil && err != iam.ErrUserNotFound {
		return nil, errors.Wrap(err, "iam.DeleteUser")
	}
	s.log.Debug("deleted user",
		zap.String("userID", req.ID),
		zap.String("username", req.Username),
		zap.Bool("deleteProjects", req.DeleteProjects))
	return &api.Void{}, nil
}
