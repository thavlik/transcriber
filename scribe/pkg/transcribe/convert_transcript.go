package transcribe

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
)

// ConvertTranscript crudely converts a transcript object into a byte slice
func ConvertTranscript(
	transcript *transcribestreamingservice.MedicalTranscript,
) string {
	s := ""
	for _, result := range transcript.Results {
		for _, alt := range result.Alternatives {
			for _, item := range alt.Items {
				if aws.StringValue(item.Type) == "punctuation" {
					continue
				}
				if item.Content == nil {
					continue
				}
				s += fmt.Sprintf("%s ", *item.Content)
			}
			s += "\n"
		}
	}
	return s
}
