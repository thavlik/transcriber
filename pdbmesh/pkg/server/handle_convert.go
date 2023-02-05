package server

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pdbmesh/pkg/convert"

	"go.uber.org/zap"
)

func (s *Server) handleConvert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			model, err := convert.Convert(
				r.Context(),
				r.Body,
			)
			if err != nil {
				return errors.Wrap(err, "failed to convert pdb to mesh")
			}
			defer model.Dispose()
			w.Header().Set("Content-Type", "model/stl")
			if _, err := io.Copy(w, model.Reader()); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
