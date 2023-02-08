package mp3

import (
	"context"
	"io"

	"github.com/hajimehoshi/go-mp3"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/scribe/pkg/source"
)

type mp3Source struct {
	ctx context.Context
	dec *mp3.Decoder
}

func NewMP3Source(
	ctx context.Context,
	r io.ReadSeeker,
) (source.Source, error) {
	dec, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, errors.Wrap(err, "mp3.NewDecoder")
	}
	return &mp3Source{
		ctx,
		dec,
	}, nil
}

func (s *mp3Source) IsStereo() (bool, error) {
	switch n := s.dec.NumChannels(); n {
	case 1:
		return false, nil
	case 2:
		return true, nil
	default:
		return false, errors.Errorf("unsupported number of channels %d", n)
	}
}

func (s *mp3Source) Context() context.Context {
	return s.ctx
}

func (s *mp3Source) SampleRate() (int64, error) {
	return int64(s.dec.SampleRate()), nil
}

func (s *mp3Source) Read(
	buf []byte,
) (int, error) {
	if err := s.ctx.Err(); err != nil {
		return 0, err
	}
	return s.dec.Read(buf)
}

func (s *mp3Source) String() string {
	return "[mp3 audio source]"
}
