package aac

import (
	"context"
	"fmt"
	"io"

	"github.com/izern/go-fdkaac/fdkaac"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AACSource struct {
	ctx        context.Context
	cancel     context.CancelFunc
	r          io.Reader
	w          io.Writer
	dec        *fdkaac.AacDecoder
	sampleRate int64
	stereo     bool
	log        *zap.Logger
}

// NewAACSource creates a new Source that reads audio from an AAC
// stream. The sampleRate parameter is the sample rate of the
// audio stream in Hertz. It is required here because the sample
// rate otherwise cannot be known until an AudioSpecificConfig (asc)
// header is read, but it's needed to create the transcription request.
// Use this Source in concert with OBS (https://obsproject.com/)
// to transcribe live audio from a microphone.
func NewAACSource(
	ctx context.Context,
	sampleRate int64,
	stereo bool,
	log *zap.Logger,
) (*AACSource, error) {
	ctx, cancel := context.WithCancel(ctx)
	r, w := io.Pipe()
	return &AACSource{
		ctx:        ctx,
		cancel:     cancel,
		r:          r,
		w:          w,
		dec:        fdkaac.NewAacDecoder(),
		sampleRate: sampleRate,
		stereo:     stereo,
		log:        log,
	}, nil
}

func (s *AACSource) InitSeqHeader(asc []byte) error {
	return s.dec.InitRaw(asc)
}

func (s *AACSource) IsStereo() (bool, error) {
	return s.stereo, nil
}

func (s *AACSource) SampleRate() (int64, error) {
	return s.sampleRate, nil
}

// Write writes a single AAC audio frame to the source.
// This method is intended to be used by the RTMP server.
func (s *AACSource) Write(frame []byte) (int, error) {
	return s.w.Write(frame)
}

func (s *AACSource) Close() error {
	s.cancel()
	return s.dec.Close()
}

func (s *AACSource) Context() context.Context {
	return s.ctx
}

func (s *AACSource) ReadAudioChunk(
	buf []byte,
) (int, error) {
	if err := s.ctx.Err(); err != nil {
		return 0, err
	}
	// read aac audio from the reader into the buffer
	n, err := s.r.Read(buf)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read aac audio")
	}
	// decode the aac audio into pcm values
	pcm, err := s.dec.Decode(buf[:n])
	if err != nil {
		s.log.Error("failed to decode aac audio", zap.Error(err))
		return 0, errors.Wrap(err, "failed to decode aac audio")
	} else if pcm == nil {
		// No audio chunk to read yet, callee should try again.
		// This should only happen if there are issues with the
		// underlying reader. The buffer should be large enough
		// to hold multiple AAC audio chunks.
		s.log.Warn("no pcm audio data available yet")
		return 0, nil
	}
	// copy the pcm audio to the output buffer
	return copy(buf, pcm), nil
}

func (s *AACSource) String() string {
	return fmt.Sprintf(
		"AACSource{sampleRate=%d, stereo=%t}",
		s.sampleRate,
		s.stereo,
	)
}
