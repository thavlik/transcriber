package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch"

	"go.uber.org/zap"
)

func (s *Server) handleImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			img := new(imgsearch.Image)
			switch r.Method {
			case http.MethodGet:
				// extract metadata from query parameters
				// this handler is provided in case a client requires
				// a direct link to an image via GET request
				if err := extractMeta(r, img); err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "extractMeta")
				}
			case http.MethodPost:
				// extract metadata from request body
				// this is the faster way of getting an image
				// when the metadata is already unmarshalled
				if r.Header.Get("Content-Type") != "application/json" {
					retCode = http.StatusUnsupportedMediaType
					return fmt.Errorf("unsupported media type %s", r.Header.Get("Content-Type"))
				}
				if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
					retCode = http.StatusBadRequest
					return errors.Wrap(err, "decode")
				}
			default:
				retCode = http.StatusMethodNotAllowed
				return fmt.Errorf("method %s not allowed", r.Method)
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
			http.Error(w, err.Error(), retCode)
		}
	}
}
