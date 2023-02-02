package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"

	"github.com/thavlik/transcriber/base/pkg/base"

	"go.uber.org/zap"
)

func (s *Server) handleImage() http.HandlerFunc {
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
			img, err := extractMeta(r)
			if err != nil {
				retCode = http.StatusBadRequest
				return err
			}
			contentLength, err := img.ContentLength()
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "parse content length")
			}
			w.Header().Set("Content-Type", img.ContentType())
			w.Header().Set("Content-Length", contentLength)
			metaHash := img.Hash()
			body, err := s.imageCache.Get(r.Context(), metaHash)
			if err == data.ErrNotCached {
				if err := s.cacheImage(
					r.Context(),
					img,
					w,
				); err != nil {
					return err
				}
				// in the event that multiple people request
				// the same uncached image at the same time,
				// we want to increment the request counter
				// both times
				s.spawn(func() {
					s.incrementRequests(metaHash)
				})
				return nil
			} else if err != nil {
				return err
			}
			defer body.Close()
			if _, err := io.Copy(w, body); err != nil {
				return errors.Wrap(err, "copy")
			}
			s.spawn(func() {
				s.incrementRequests(metaHash)
			})
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}
