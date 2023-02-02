package raw

import (
	"context"
	"io"

	"github.com/thavlik/transcriber/transcriber/pkg/source"
)

type rawSource struct {
	ctx        context.Context
	cancel     context.CancelFunc
	sampleRate int64
	isStereo   bool
	r          io.Reader
}

func NewRawSource(
	ctx context.Context,
	sampleRate int64,
	isStereo bool,
	r io.Reader,
) source.Source {
	ctx, cancel := context.WithCancel(ctx)
	return &rawSource{
		ctx:        ctx,
		cancel:     cancel,
		sampleRate: sampleRate,
		isStereo:   isStereo,
		r:          r,
	}
}

func (s *rawSource) Read(buf []byte) (int, error) {
	if s.ctx.Err() != nil {
		return 0, s.ctx.Err()
	}
	n, err := s.r.Read(buf)
	if err != nil {
		s.cancel()
		return 0, err
	}
	return n, nil
}

func (s *rawSource) IsStereo() (bool, error) {
	return s.isStereo, nil
}

func (s *rawSource) SampleRate() (int64, error) {
	return s.sampleRate, nil
}

func (s *rawSource) Context() context.Context {
	return s.ctx
}
