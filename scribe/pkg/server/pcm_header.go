package server

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type pcmHeader struct {
	sampleRate int64
	isStereo   bool
}

var errInvalidHeader = errors.New("invalid content type")

func parsePCMHeader(contentType string) (hdr *pcmHeader, err error) {
	parts := strings.Split(contentType, ";")
	if len(parts) != 3 {
		return nil, errInvalidHeader
	}
	if parts[0] != "audio/pcm" {
		return nil, errInvalidHeader
	}
	parts = parts[1:]
	hdr = new(pcmHeader)
	for _, part := range parts {
		parts := strings.Split(part, "=")
		if len(parts) != 2 {
			return nil, errInvalidHeader
		}
		switch parts[0] {
		case "bits":
			switch parts[1] {
			case "16":
				// do nothing, default is 16 bits
			default:
				return nil, errors.Errorf("invalid bit count '%s', must be '16'", parts[1])
			}
		case "encoding":
			switch parts[1] {
			case "signed-integer":
				// do nothing, default is signed-integer
			default:
				return nil, errors.Errorf("invalid encoding '%s', must be 'signed-integer'", parts[1])
			}
		case "big-endian":
			switch parts[1] {
			case "true":
				return nil, errors.New("big-endian not supported")
			case "false":
				// do nothing, default is little-endian
			default:
				return nil, errors.Errorf("invalid big-endian value '%s', must be 'true' or 'false'", parts[1])
			}
		case "rate":
			hdr.sampleRate, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse sample rate")
			}
		case "channels":
			switch parts[1] {
			case "1":
				// do nothing, default is mono
			case "2":
				hdr.isStereo = true
			default:
				return nil, errors.Errorf("invalid channel count '%s', must be '1' or '2'", parts[1])
			}
		default:
			return nil, errors.Errorf("invalid option %s", parts[0])
		}
	}
	return
}
