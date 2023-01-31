package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"

	"github.com/thavlik/transcriber/base/pkg/base"

	"go.uber.org/zap"
)

func (s *server) handleImage() http.HandlerFunc {
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
			body, err := s.imageCache.Get(r.Context(), img.Hash())
			if err == data.ErrNotCached {
				if err := cacheImage(
					context.TODO(), // don't use the request context here
					img,
					s.imageCache,
				); err != nil {
					s.log.Error("failed to cache image", zap.Error(err))
					return err
				}
				body, err = s.imageCache.Get(r.Context(), img.Hash())
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			defer body.Close()
			if _, err := io.Copy(w, body); err != nil {
				return errors.Wrap(err, "copy")
			}
			return nil
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			w.WriteHeader(retCode)
			w.Write([]byte(err.Error()))
		}
	}
}

func extractMeta(r *http.Request) (*search.Image, error) {
	input := r.URL.Query().Get("i")
	if input == "" {
		return nil, errors.New("missing query parameter 'i'")
	}
	unescaped, err := url.QueryUnescape(input)
	if err != nil {
		return nil, errors.Wrap(err, "url.QueryUnescape")
	}
	img := new(search.Image)
	if err := json.Unmarshal(
		[]byte(unescaped),
		&img,
	); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return img, nil
}

func cacheImage(
	ctx context.Context,
	img *search.Image,
	imageCache *cache.ImageCache,
) error {
	req, err := http.NewRequest(
		"GET",
		img.ContentURL,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "newrequest")
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != img.ContentType() {
		// sanity check to ensure we're not caching the wrong thing
		// if this fails, it's probably an image with a rare mimetype
		// in any case, we probably don't want to show it to the user
		return fmt.Errorf(
			"content type mismatch: expected %s, got %s",
			img.ContentType(),
			resp.Header.Get("Content-Type"),
		)
	}
	if err := imageCache.Set(
		ctx,
		img,
		resp.Body,
	); err != nil {
		return errors.Wrap(err, "imagecache.Set")
	}
	return nil
}
