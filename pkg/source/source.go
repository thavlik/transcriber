package source

type Source interface {
	// SampleRate returns the sample rate of the audio source in Hertz.
	SampleRate() (int64, error)

	// Encoding returns the encoding of the audio source. This is usually "pcm".
	Encoding() (string, error)

	// ReadAudioChunk reads a chunk of audio from the source.
	ReadAudioChunk(buf []byte) (n int, err error)
}
