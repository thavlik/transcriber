package base

import (
	"context"
	"time"
)

type DownloadProgress struct {
	Total   int64         // bytes
	Rate    float64       // bytes per second
	Elapsed time.Duration //
}

func Progress(ctx context.Context, onProgress chan<- struct{}) {
	if onProgress != nil {
		select {
		case <-ctx.Done():
			return
		case onProgress <- struct{}{}:
		}
	}
}

func ProgressDownload(
	ctx context.Context,
	onProgress chan<- *DownloadProgress,
) {
	if onProgress != nil {
		select {
		case <-ctx.Done():
			return
		case onProgress <- nil:
			// sending a nil value to this channel is
			// the same as sending the empty struct
			// to the untyped progress channel
		}
	}
}
