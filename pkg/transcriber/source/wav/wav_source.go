package file

import (
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/pkg/transcriber/source"
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

func (s *wavSource) SampleRate() (int64, error) {
	return int64(s.dec.SampleRate), nil
}

func (s *wavSource) Encoding() (string, error) {
	// all wav files are PCM
	return "pcm", nil
}

func (s *wavSource) ReadAudioChunk(
	buf []byte,
) (int, error) {
	n, err := s.dec.PCMBuffer(s.buf)
	if err != nil {
		return 0, errors.Wrap(err, "dec.PCMBuffer")
	} else if n == 0 {
		return 0, io.EOF
	} else if s.buf.SourceBitDepth != 16 {
		return 0, errors.Errorf("unsupported bit depth %d", s.buf.SourceBitDepth)
	}
	// Convert to PCM signed 16 bit little endian as per:
	// https://docs.aws.amazon.com/transcribe/latest/dg/how-input.html
	// wav file PCM data is already little endian
	for i := 0; i < n; i++ {
		buf[2*i] = byte(s.buf.Data[i] & 0xff)
		buf[2*i+1] = byte((s.buf.Data[i] >> 8) & 0xff)
	}
	return n * 2, nil
}

func (s *wavSource) String() string {
	return s.dec.String()
}
