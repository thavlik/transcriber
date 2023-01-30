package wav

import (
	"context"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/transcriber/pkg/source"
)

type wavSource struct {
	dec *wav.Decoder
	buf *audio.IntBuffer
}

func NewWavSource(
	r io.ReadSeeker,
) (source.Source, error) {
	dec := wav.NewDecoder(r)
	if !dec.IsValidFile() {
		if err := dec.Err(); err != nil {
			return nil, errors.Wrap(err, "dec.ReadInfo")
		}
		return nil, errors.New("invalid wav file")
	}
	if err := dec.FwdToPCM(); err != nil {
		return nil, errors.Wrap(err, "dec.FwdToPCM")
	}
	if dec.BitDepth != 16 {
		// TODO: support other bit depths
		return nil, errors.Errorf("unsupported bit depth %d", dec.BitDepth)
	}
	return &wavSource{
		dec,
		&audio.IntBuffer{
			Data: make([]int, 15000), // max chunk size is 32kb, so 15000 16 bit samples = 30kb
		},
	}, nil
}

func (s *wavSource) IsStereo() (bool, error) {
	switch s.dec.NumChans {
	case 1:
		return false, nil
	case 2:
		return true, nil
	default:
		return false, errors.Errorf("unsupported number of channels %d", s.dec.NumChans)
	}
}

func (s *wavSource) Context() context.Context {
	return context.Background()
}

func (s *wavSource) SampleRate() (int64, error) {
	return int64(s.dec.SampleRate), nil
}

func (s *wavSource) ReadAudioChunk(
	buf []byte,
) (int, error) {
	n, err := s.dec.PCMBuffer(s.buf)
	if err != nil {
		return 0, errors.Wrap(err, "dec.PCMBuffer")
	} else if n == 0 {
		return 0, io.EOF
	}
	// Convert to PCM signed 16 bit little endian as per:
	// https://docs.aws.amazon.com/transcribe/latest/dg/how-input.html
	// wav file PCM data is already little endian
	for i := 0; i < n; i++ {
		v := s.buf.Data[i]
		buf[2*i] = byte(v & 0xff)
		buf[2*i+1] = byte((v >> 8) & 0xff)
	}
	return n * 2, nil
}

func (s *wavSource) String() string {
	return s.dec.String()
}
