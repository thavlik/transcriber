package server

import (
	"io"
	"net/http"
)

func (s *Server) handlePublish() http.HandlerFunc {
	return s.handler(
		http.MethodPost,
		func(w http.ResponseWriter, r *http.Request) error {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			s.spawn(func() {
				s.publish(body)
			})
			return nil
		},
	)
}
