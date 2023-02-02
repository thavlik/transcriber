package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"

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
			var historyDone chan error
			if req.UserID != "" {
				historyDone = make(chan error, 1)
				s.spawn(func() {
					historyDone <- s.history.Push(
						s.ctx,
						&history.Search{
							ID:        uuid.New().String(),
							Query:     req.Query,
							UserID:    req.UserID,
							Timestamp: time.Now(),
						},
					)
				})
			}
			images, err := search.Search(
				r.Context(),
				req.Query,
				s.endpoint,
				s.apiKey,
				10,
				0,
			)
			if err != nil {
				return errors.Wrap(err, "search failed")
			}
			if historyDone != nil {
				select {
				case <-r.Context().Done():
					return r.Context().Err()
				case err := <-historyDone:
					if err != nil {
						return errors.Wrap(err, "failed to push search to history")
					}
				}
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
