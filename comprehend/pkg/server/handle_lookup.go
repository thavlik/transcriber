package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
	"go.uber.org/zap"
)

func (s *Server) handleLookup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing id")
			}
			entity, err := s.entityCache.Lookup(
				r.Context(),
				id,
			)
			if err != nil {
				if err == entitycache.ErrNotFound {
					retCode = http.StatusNotFound
				}
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(entity)
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
