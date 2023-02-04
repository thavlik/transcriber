package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/search/adapter"

	"go.uber.org/zap"
)

func (s *Server) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			var req struct {
				Query  string `json:"query"`
				UserID string `json:"userID,omitempty"`
			}
			switch r.Method {
			case http.MethodOptions:
				base.AddPreflightHeaders(w)
				return nil
			case http.MethodGet:
				req.Query = r.URL.Query().Get("q")
			case http.MethodPost:
				if r.Header.Get("Content-Type") != "application/json" {
					retCode = http.StatusUnsupportedMediaType
					return fmt.Errorf("unsupported media type %s", r.Header.Get("Content-Type"))
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "failed to decode request body")
				}
			default:
				retCode = http.StatusMethodNotAllowed
				return errors.New("bad method")
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if req.Query == "" {
				retCode = http.StatusBadRequest
				return fmt.Errorf("missing query")
			}
			start := time.Now()
			service := adapter.Bing
			images, err := adapter.Search(
				r.Context(),
				service,
				req.Query,
				s.endpoint,
				s.apiKey,
				10,
				0,
			)
			if err != nil {
				return errors.Wrap(err, "search failed")
			}
			s.log.Debug(
				"searched for images",
				base.Elapsed(start),
				zap.String("service", string(service)),
				zap.Int("count", len(images)),
			)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(images); err != nil {
				return err
			}
			if req.UserID != "" {
				s.spawn(func() {
					// Only count searches for logged in users,
					// don't block the response, and only count
					// it if everything is successful. This is
					// the most generous policy for the user.
					s.pushSearchHistory(req.Query, req.UserID)
				})
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
