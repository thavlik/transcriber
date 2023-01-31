package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/imgsearch/pkg/search"
)

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
		img,
	); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return img, nil
}
