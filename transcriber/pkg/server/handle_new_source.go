package server

import (
	"net/http"

	"github.com/thavlik/transcriber/transcriber/pkg/source/raw"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Server) handleNewSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() error {
			if r.Method != http.MethodPost {
				retCode = http.StatusMethodNotAllowed
				return errors.New("method not allowed")
			}
			// parse the pcm header
			// the header should be in the format:
			// audio/pcm;rate=44100;channels=1
			// https://voysis.readme.io/docs/audio-guidelines
			hdr, err := parsePCMHeader(r.Header.Get("Content-Type"))
			if err != nil {
				retCode = http.StatusBadRequest
				return errors.Wrap(err, "invalid content type")
			}
			defer r.Body.Close()
			source := raw.NewRawSource(
				r.Context(),
				hdr.sampleRate,
				hdr.isStereo,
				r.Body,
			)
			// try and assign a new audio source
			select {
			case <-r.Context().Done():
				return r.Context().Err()
			case s.newSource <- source:
				// assigned new source & created new job
			}
			// wait for job to finish or the request to cancel
			select {
			case <-r.Context().Done():
				return r.Context().Err()
			case <-source.Context().Done():
				return errors.Wrap(source.Context().Err(), "source error")
			}
		}(); err != nil {
			s.log.Error(r.RequestURI, zap.Error(err))
			http.Error(w, err.Error(), retCode)
		}
	}
}
