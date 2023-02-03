package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"go.uber.org/zap"
)

func (s *Server) handleDetect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			defer r.Body.Close()
			if r.Header.Get("Content-Type") != "application/json" {
				retCode = http.StatusUnsupportedMediaType
				return errors.New("unsupported media type")
			}
			var req struct {
				Text         string             `json:"text"`
				Filter       *comprehend.Filter `json:"filter,omitempty"`
				Model        string             `json:"model,omitempty"`
				LanguageCode string             `json:"languageCode,omitempty"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				retCode = http.StatusBadRequest
				return err
			}
			if req.Text == "" {
				retCode = http.StatusBadRequest
				return errors.New("missing text")
			}
			if req.LanguageCode == "" {
				req.LanguageCode = "en"
			}
			var entities []*comprehend.Entity
			switch req.Model {
			case "":
				fallthrough
			case "amazon-comprehend":
				if entities, err = comprehend.Comprehend(
					r.Context(),
					req.Text,
					req.LanguageCode,
					req.Filter,
				); err != nil {
					return err
				}
			case "amazon-comprehend-medical":
				if entities, err = comprehend.ComprehendMedical(
					r.Context(),
					req.Text,
					req.Filter,
				); err != nil {
					return err
				}
			default:
				retCode = http.StatusBadRequest
				return errors.Errorf("invalid model '%s'", req.Model)
			}
			s.spawn(func() {
				// cache the entity hashes so we can restrict what
				// inputs can be used to services like OpenAI/GPT-3
				if err := s.entityCache.BulkCache(
					r.Context(),
					entities,
				); err != nil {
					s.log.Error("failed to cache entities", zap.Error(err))
				}
			})
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(entities)
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
