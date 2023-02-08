package batch

import (
	"fmt"
	"path/filepath"
)

type BatchTranscribeSource struct {
	Bucket string
	Key    string
}

// Format returns the file extension of the source file.
func (b *BatchTranscribeSource) Format() string {
	ext := filepath.Ext(b.Key)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

// S3Uri returns the S3 URI of the source file.
func (b *BatchTranscribeSource) S3Uri() string {
	return fmt.Sprintf("s3://%s/%s", b.Bucket, b.Key)
}
