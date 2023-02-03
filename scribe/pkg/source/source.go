package source

import (
	"context"
)

// Source is a Reader-derived interface for reading PCM audio data from a source.
// The audio data is expected to be signed 16 bit little endian.
// If the source has a different bit depth, it should be converted
// to 16 bit before being returned.
type Source interface {
	// Context returns the context for the source.
	// Cancelling the provided context will terminate decoding
	// and close the source file/stream/etc.
	Context() context.Context

	// SampleRate returns the sample rate of the audio source in Hertz.
	SampleRate() (int64, error)

	// IsStereo returns true if the audio source is stereo.
	// This should always be true for RTMP streams but may
	// or may not be true for other sources, e.g. wav files.
	// Amazon Transcribe requires either mono or stereo, so
	// if the source is 5.1 surround sound or similar, this
	// should return an error (unless you want to downmix)
	IsStereo() (bool, error)

	// Read reads the next PCM chunk from the source.
	// The chunk is expected to be signed 16 bit little endian.
	Read(buf []byte) (n int, err error)
}
