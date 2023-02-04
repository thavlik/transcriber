package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/disease"
	"github.com/thavlik/transcriber/define/pkg/diseasecache"
	"go.uber.org/zap"
)

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
				retCode = http.StatusLoopDetected
				return errors.New("missing query")
			}
			start := time.Now()
			query = strings.TrimSpace(query)
			if isDisease, err := s.diseaseCache.IsDisease(
				r.Context(),
				query,
			); err == nil {
				s.log.Debug("disease was cached",
					zap.String("query", query),
					zap.Bool("isDisease", isDisease),
					base.Elapsed(start))
				// use the cached value
				w.Header().Set("Content-Type", "application/json")
				return json.NewEncoder(w).Encode(map[string]interface{}{
					"isDisease": isDisease,
				})
			} else if err != diseasecache.ErrNotFound {
				return errors.Wrap(err, "disease cache failed")
			}
			start = time.Now()
			isDisease, err := disease.IsDisease(
				r.Context(),
				s.gpt3,
				query,
			)
			if err != nil {
				return errors.Wrap(err, "disease.IsDisease")
			}
			s.log.Debug("queried gpt3 for disease",
				zap.String("query", query),
				zap.Bool("isDisease", isDisease),
				base.Elapsed(start))
			s.spawn(func() {
				if err := s.diseaseCache.Set(
					s.ctx,
					query,
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
