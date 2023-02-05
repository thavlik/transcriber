package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func (s *Server) getPDB(
	ctx context.Context,
	id string,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/pdb?id=%s", s.pharmaSeerOpts.Endpoint, id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{
		Timeout: s.pharmaSeerOpts.Timeout,
	}).Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"failed to download pdb for %s: status code %d",
			id,
			resp.StatusCode,
		)
	}
	return resp.Body, nil
}
