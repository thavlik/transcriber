package transcriber

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
)

func PrintTranscripts(
	ctx context.Context,
	transcripts <-chan *transcribestreamingservice.MedicalTranscript,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case transcript := <-transcripts:
			for _, result := range transcript.Results {
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
		}
	}
}
