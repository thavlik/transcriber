package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thavlik/transcriber/imgsearch/pkg/search"

	"go.uber.org/zap"
)

func (s *server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodGet {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			images, err := search.Search(
				r.Context(),
				r.URL.Query().Get("q"),
				s.endpoint,
				s.apiKey,
				10,
				0,
			)
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(images)
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
