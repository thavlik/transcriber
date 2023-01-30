package base

import (
	"strings"

	"github.com/pkg/errors"
)

func ExtractChannelID(input string) (string, error) {
	input = strings.ReplaceAll(input, "https://", "")
	input = strings.ReplaceAll(input, "http://", "")
	input = strings.ReplaceAll(input, "www.", "")
	input = strings.ReplaceAll(input, "m.", "")
	if i := strings.Index(input, "/"); i != -1 {
		input = input[i+1:]
	}
	if !strings.HasPrefix(input, "@") {
		return "", errors.New("username is missing @ prefix")
	}
	return input, nil
}
