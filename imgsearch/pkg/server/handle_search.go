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
			// TODO: check iam
			userID := "test"
			h := &history.Search{
				ID:        uuid.New().String(),
				Query:     r.URL.Query().Get("q"),
				UserID:    userID,
				Timestamp: time.Now(),
			}
			if h.Query == "" {
				retCode = http.StatusBadRequest
				return fmt.Errorf("query parameter 'q' is required")
			}
			historyDone := make(chan error, 1)
			s.spawn(func() {
				historyDone <- s.history.Push(
					s.ctx,
					h,
				)
			})
			images, err := search.Search(
				r.Context(),
				h.Query,
				s.endpoint,
				s.apiKey,
				10,
				0,
			)
			if err != nil {
				return errors.Wrap(err, "search failed")
			}
			select {
			case <-r.Context().Done():
				return r.Context().Err()
			case err := <-historyDone:
				if err != nil {
					return errors.Wrap(err, "failed to push search to history")
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
