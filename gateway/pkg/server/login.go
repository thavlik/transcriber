package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam/api"
	"go.uber.org/zap"
)

func (s *Server) Login(
	ctx context.Context,
	req api.LoginRequest,
) (*api.LoginResponse, error) {
	if req.Username == "" {
		return nil, errors.New("missing username")
	}
	if req.Password == "" {
		return nil, errors.New("missing password")
	}
	accessToken, err := s.iam.Login(
		ctx,
		req.Username,
		req.Password,
	)
	if err != nil {
		return nil, errors.Wrap(err, "iam.Login")
	}
	s.log.Debug("login success", zap.String("username", req.Username))
	return &api.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
