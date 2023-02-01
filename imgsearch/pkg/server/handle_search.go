package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"

	"go.uber.org/zap"
)

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			switch r.Method {
			case http.MethodOptions:
				base.AddPreflightHeaders(w)
				return nil
			case http.MethodGet:
				break
			default:
				retCode = http.StatusMethodNotAllowed
				return fmt.Errorf("method not allowed")
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
