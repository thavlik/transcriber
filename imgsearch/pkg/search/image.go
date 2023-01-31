package search

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type Image struct {
	ContentURL     string `json:"contentURL"`
	ContentSize    string `json:"contentSize"`
	ThumbnailURL   string `json:"thumbnailURL"`
	HostPageURL    string `json:"hostPageURL"`
	EncodingFormat string `json:"encodingFormat"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	AccentColor    string `json:"accentColor"`
}

func (i *Image) Hash() string {
	h := md5.New()
	// note: md5.Write cannot fail so we don't need to check the error
	h.Write([]byte(i.ContentURL))
	h.Write([]byte(i.ContentSize))
	h.Write([]byte(i.ThumbnailURL))
	h.Write([]byte(i.HostPageURL))
	h.Write([]byte(i.EncodingFormat))
	h.Write([]byte(fmt.Sprintf("%d\n%d\n", i.Width, i.Height)))
	raw := h.Sum(nil) // 16 bytes
	return hex.EncodeToString(raw[:])
}

func (i *Image) ContentLength() (string, error) {
	parts := strings.Split(i.ContentSize, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid content size: %s", i.ContentSize)
	} else if parts[1] != "B" {
		return "", fmt.Errorf("invalid content size: %s", i.ContentSize)
	}
	return parts[0], nil
}

func (i *Image) ContentType() string {
	switch i.EncodingFormat {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "bmp":
		return "image/bmp"
	case "tiff":
		return "image/tiff"
	case "svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
