package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/storage"
	"go.uber.org/zap"
)

func (s *Server) handleDefine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			// TODO: input should be a UUID to an entity
			// so only definitions of recognized entities
			// are possible. Currently, this exposes the
			// full GPT-3 API to the world.
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			query := r.URL.Query().Get("q")
			if query == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing query")
			}
			query = strings.TrimSpace(query)
			n := 1
			var temp float32 = 0.7
			var topP float32 = 1.0
			maxLength := 256
			timestamp := time.Now()
			resp, err := s.gpt3.Completion(
				r.Context(),
				gpt3.CompletionRequest{
					Prompt:           []string{query},
					Temperature:      &temp,
					MaxTokens:        &maxLength,
					TopP:             &topP,
					N:                &n,
					FrequencyPenalty: 0.0,
					PresencePenalty:  0.0,
				},
			)
			if err != nil {
				return errors.Wrap(err, "gpt3")
			}
			output := strings.TrimSpace(resp.Choices[0].Text)
			s.spawn(func() {
				if err := s.storage.Insert(
					s.ctx,
					&storage.Definition{
						ID:        uuid.New().String(),
						Input:     query,
						Output:    output,
						Timestamp: timestamp,
					},
				); err != nil {
					s.log.Error("failed to save definition", zap.Error(err))
				}
			})
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(map[string]interface{}{
				"text": output,
			})
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
