package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/scribe/pkg/source"
)

func (s *Server) pushAudioSource(
	ctx context.Context,
	src source.Source,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.scribe.Endpoint+"/source",
		src,
	)
	if err != nil {
		return err
	}
	if s.scribe.BasicAuth.Username != "" {
		req.SetBasicAuth(
			s.scribe.BasicAuth.Username,
			s.scribe.BasicAuth.Password,
		)
	}
	rate, err := src.SampleRate()
	if err != nil {
		return err
	}
	isStereo, err := src.IsStereo()
	if err != nil {
		return err
	}
	var channels int
	if isStereo {
		channels = 2
	} else {
		channels = 1
	}
	req.Header.Set(
		"Content-Type",
		fmt.Sprintf(
			"audio/pcm;rate=%d;channels=%d",
			rate,
			channels,
		),
	)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return errors.Errorf("status code %d: %s", resp.StatusCode, string(body))
}
