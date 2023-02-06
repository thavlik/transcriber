package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (s *Server) handleDefine() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query().Get("q")
			if query == "" {
				return errors.New("missing query parameter 'q'")
			}
			// create an appropriate completion prompt for ChatGPT
			s.log.Debug("defining term", zap.String("query", query))
			query = url.QueryEscape(fmt.Sprintf(`define "%s"`, query))
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodGet,
				s.define.Endpoint+"/completion?q="+query,
				nil,
			)
			if err != nil {
				return err
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := (&http.Client{
				Timeout: s.define.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("define service status code %d: %s", resp.StatusCode, string(body))
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Server) handleIsDisease() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query().Get("q")
			if query == "" {
				return errors.New("missing query parameter 'q'")
			}
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodGet,
				s.define.Endpoint+"/disease?q="+url.QueryEscape(query),
				nil,
			)
			if err != nil {
				return err
			}
			resp, err := (&http.Client{
				Timeout: s.define.Timeout,
			}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("define service status code %d: %s", resp.StatusCode, string(body))
			}
			w.Header().Set("Content-Type", "application/json")
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		},
	)
}
