package base

func Duplicate(buf []byte) []byte {
	dup := make([]byte, len(buf))
	copy(dup, buf)
	return dup
}
