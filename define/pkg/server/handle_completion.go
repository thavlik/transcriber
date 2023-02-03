package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) handleCompletion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			query := r.URL.Query().Get("q")
			if query == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing query")
			}
			unescaped, err := url.QueryUnescape(query)
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "unescaping query")
			}
			unescaped = strings.TrimSpace(unescaped)
			client := gpt3.NewClient(
				s.openAISecretKey,
				gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
			)
			n := 1
			var temp float32 = 0.7
			var topP float32 = 1.0
			maxLength := 256
			resp, err := client.Completion(
				r.Context(),
				gpt3.CompletionRequest{
					Prompt:           []string{unescaped},
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
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(map[string]interface{}{
				"text": strings.TrimSpace(resp.Choices[0].Text),
			})
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
