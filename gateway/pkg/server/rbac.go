package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrMissingAccessToken = errors.New("missing AccessToken in header")
)

func (s *Server) rbac(
	ctx context.Context,
	r *http.Request,
	permissions []string,
) (string, error) {
	if s.iam == nil {
		// RBAC is disabled. This is only suitable for development.
		return "", nil
	}
	accessToken := r.Header.Get("AccessToken")
	if accessToken == "" {
		return "", ErrMissingAccessToken
	}
	userID, err := s.iam.Authorize(
		ctx,
		accessToken,
		permissions,
	)
	if err != nil {
		return "", errors.Wrap(err, "iam")
	}
	return userID, nil
}
