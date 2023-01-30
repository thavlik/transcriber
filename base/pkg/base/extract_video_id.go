package base

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var ErrNoVideoID = errors.New("url query is missing video id")

func ExtractVideoID(input string) (string, error) {
	if strings.Contains(input, ".") {
		u, err := url.Parse(input)
		if err != nil {
			return "", errors.Wrap(err, "url.Parse")
		}
		v := u.Query().Get("v")
		if v == "" {
			return "", ErrNoVideoID
		}
		return v, nil
	}
	if n := len(input); n != 11 {
		return "", fmt.Errorf("video ID '%s' has invalid length %d", input, n)
	}
	return input, nil
}
