package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"go.uber.org/zap"
)

type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

func (s *Server) broadcastMessage(
	ctx context.Context,
	ty string,
	payload interface{},
) {
	body, err := json.Marshal(&wsMessage{
		Type:    ty,
		Payload: payload,
	})
	if err != nil {
		panic(err)
	}
	s.broadcast(ctx, body)
}

func (s *Server) broadcast(
	ctx context.Context,
	body []byte,
) {
	if err := func() error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.broadcaster.Endpoint+"/publish",
			bytes.NewReader(body),
		)
		if err != nil {
			return err
		}
		resp, err := (&http.Client{
			Timeout: s.broadcaster.Timeout,
		}).Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return errors.Errorf(
				"status code %d: %s",
				resp.StatusCode,
				string(body),
			)
		}
		return nil
	}(); err != nil {
		s.log.Error("failed to publish through broadcaster", zap.Error(err))
	}
}
