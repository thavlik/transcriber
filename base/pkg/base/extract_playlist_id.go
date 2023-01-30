package base

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var ErrNoPlaylistID = errors.New("url query is missing playlist id")

func ExtractPlaylistID(input string) (string, error) {
	if strings.Contains(input, ".") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.Wrap(err, "url.Parse")
		}
		v := u.Query().Get("list")
		if v == "" {
			return "", ErrNoPlaylistID
		}
	}
	// further verification may be more difficult
	// than simply reaching out to youtube and
	// seeing what we get
	return input, nil
}
