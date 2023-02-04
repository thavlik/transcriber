package transcribe

import (
	"context"
	"fmt"
)

func PrintTranscripts(
	ctx context.Context,
	transcripts <-chan *Transcript,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case transcript, ok := <-transcripts:
			if !ok {
				return
			}
			fmt.Println(transcript.String())
		}
	}
}
