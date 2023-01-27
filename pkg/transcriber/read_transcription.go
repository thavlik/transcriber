package transcriber

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
)

func readTranscription(
	ctx context.Context,
	events <-chan transcribestreamingservice.TranscriptResultStreamEvent,
) error {
	for ev := range events {
		if ev == nil {
			continue
		}
		if e, ok := ev.(*transcribestreamingservice.TranscriptEvent); ok {
			for _, result := range e.Transcript.Results {
				for i, alt := range result.Alternatives {
					fmt.Printf("%d: ", i)
					for _, item := range alt.Items {
						if aws.StringValue(item.Type) == "punctuation" {
							continue
						}
						if item.Content == nil {
							continue
						}
						fmt.Printf("%s ", *item.Content)
					}
					fmt.Printf("\n")
				}
			}
		} else {
			fmt.Printf("unrecognized event: %T\n", ev)
		}
	}
	return nil
}
