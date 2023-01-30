package source

import "context"

// Source is an interface for reading PCM audio data from a source.
// The audio data is expected to be signed 16 bit little endian.
type Source interface {
	Context() context.Context

	// SampleRate returns the sample rate of the audio source in Hertz.
	SampleRate() (int64, error)

	IsStereo() (bool, error)

	// ReadAudioChunk reads a chunk of audio from the source.
	// The chunk is expected to be signed 16 bit little endian.
	ReadAudioChunk(buf []byte) (n int, err error)
}
