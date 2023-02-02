package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (s *Server) handleImage() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			input := r.URL.Query().Get("i")
			if input == "" {
				return errors.New("missing query parameter 'i'")
			}
			body, err := base64.RawURLEncoding.DecodeString(input)
			if err != nil {
				return err
			}
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodPost,
				s.imgSearch.Endpoint+"/img/view",
				bytes.NewReader(body),
			)
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := (&http.Client{
				Timeout: s.imgSearch.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("imgsearch status code %d: %s", resp.StatusCode, string(body))
			}
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			if l := resp.Header.Get("Content-Length"); l != "" {
				w.Header().Set("Content-Length", l)
			}
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Server) handleImageSearch() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query().Get("q")
			if query == "" {
				return fmt.Errorf("query parameter 'q' is required")
			}
			body, err := json.Marshal(map[string]interface{}{
				"query":  query,
				"userID": userID,
			})
			if err != nil {
				return err
			}
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodPost,
				s.imgSearch.Endpoint+"/img/search",
				bytes.NewReader(body),
			)
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := (&http.Client{
				Timeout: s.imgSearch.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("imgsearch status code %d: %s", resp.StatusCode, string(body))
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		},
	)
}
