package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
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
			unescaped, err := url.QueryUnescape(query)
			if err != nil {
				return err
			}
			// create an appropriate completion prompt for ChatGPT
			query = fmt.Sprintf(`define "%s"`, unescaped)
			req, err := http.NewRequestWithContext(
				r.Context(),
				http.MethodGet,
				s.define.Endpoint+"/completion?q="+url.QueryEscape(query),
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
