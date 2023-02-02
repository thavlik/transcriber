package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

func extractMeta(r *http.Request, img *search.Image) error {
	input := r.URL.Query().Get("i")
	if input == "" {
		return errors.New("missing query parameter 'i'")
	}
	body, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return errors.Wrap(err, "url.QueryUnescape")
	}
	if err := json.Unmarshal(
		[]byte(body),
		img,
	); err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}
	return nil
}
