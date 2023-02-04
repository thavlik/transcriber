package transcribe

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/transcribestreamingservice"
)

func convertTranscript(
	input *transcribestreamingservice.MedicalTranscript,
) *Transcript {
	results := make([]*TranscriptionResult, len(input.Results))
	for i, result := range input.Results {
		alternatives := make([]string, len(result.Alternatives))
		for i, alt := range result.Alternatives {
			alternatives[i] = aws.StringValue(alt.Transcript)
		}
		results[i] = &TranscriptionResult{
			StartTime:    aws.Float64Value(result.StartTime),
			EndTime:      aws.Float64Value(result.EndTime),
			IsPartial:    aws.BoolValue(result.IsPartial),
			Alternatives: alternatives,
		}
	}
	return &Transcript{
		Results: results,
	}
}
