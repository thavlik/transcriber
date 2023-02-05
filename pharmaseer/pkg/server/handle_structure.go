package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) handleStructure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			id := r.URL.Query().Get("id")
			if id == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing query parameter 'id'")
			}
			w.Header().Set("Content-Type", "image/svg+xml")
			if err := s.svgCache.Get(
				r.Context(),
				id,
				w,
			); err == api.ErrNotCached {
				// get the thumbnail directly from drugbank
				if err := downloadDrugSVG(
					r.Context(),
					id,
					s.svgCache,
					w,
				); err != nil {
					return errors.Wrap(err, "downloadDrugSVG")
				}
				return nil
			} else if err != nil {
				return errors.Wrap(err, "svgCache.Get")
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
