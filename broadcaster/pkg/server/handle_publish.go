package server

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) handlePublish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return fmt.Errorf("method not allowed")
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			s.spawn(func() {
				s.publish(body)
			})
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
