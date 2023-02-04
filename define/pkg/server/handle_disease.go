package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/define/pkg/disease"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"go.uber.org/zap"
)

func IsDiseaseQuery(input string) string {
	return fmt.Sprintf(
		"Yes or no, is the term \"%s\" a kind of disease?",
		input,
	)
}

func (s *Server) handleDisease() http.HandlerFunc {
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
			input, err := url.QueryUnescape(query)
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "unescaping query")
			}
			input = strings.TrimSpace(input)
			if isDisease, err := s.diseaseCache.IsDisease(
				r.Context(),
				input,
			); err == nil {
				// use the cached value
				w.Header().Set("Content-Type", "application/json")
				return json.NewEncoder(w).Encode(map[string]interface{}{
					"isDisease": isDisease,
				})
			} else if err != diseasecache.ErrNotFound {
				return errors.Wrap(err, "disease cache failed")
			}
			isDisease, err := disease.IsDisease(
				r.Context(),
				s.gpt3,
				input,
			)
			if err != nil {
				return errors.Wrap(err, "disease.IsDisease")
			}
			s.spawn(func() {
				if err := s.diseaseCache.Set(
					r.Context(),
					input,
					isDisease,
				); err != nil {
					s.log.Error("failed to set disease cache", zap.Error(err))
				}
			})
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(map[string]interface{}{
				"isDisease": isDisease,
			})
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
