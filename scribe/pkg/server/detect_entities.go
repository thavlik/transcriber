package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
)

func (s *Server) detectEntities(
	ctx context.Context,
	model string,
	text string,
) ([]*comprehend.Entity, error) {
	body, err := json.Marshal(map[string]interface{}{
		"model":  model,
		"text":   text,
		"filter": s.filter,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.comprehend.Endpoint+"/detect",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{
		Timeout: s.comprehend.Timeout,
	}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("comprehend error %d: %s", resp.StatusCode, body)
	}
	var entities []*comprehend.Entity
	if err := json.NewDecoder(resp.Body).Decode(&entities); err != nil {
		return nil, err
	}
	return entities, nil
}
