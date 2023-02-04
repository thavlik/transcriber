package transcribe

import (
	"fmt"
)

type TranscriptItem struct {
	Content   string  `json:"content"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
}

type TranscriptionResult struct {
	Alternatives []string `json:"alternatives"`
	StartTime    float64  `json:"startTime"`
	EndTime      float64  `json:"endTime"`
	IsPartial    bool     `json:"isPartial"`
}

type Transcript struct {
	Results []*TranscriptionResult `json:"results"`
}

func (t *Transcript) Flatten() string {
	s := ""
	for _, result := range t.Results {
		for _, alt := range result.Alternatives {
			s += fmt.Sprintf("%s ", alt)
		}
	}
	return s
}

func (t *Transcript) String() string {
	s := "{\n"
	for _, result := range t.Results {
		for i, alt := range result.Alternatives {
			s += fmt.Sprintf("  Alt %d: %s\n", i, alt)
		}
	}
	return s + "}"
}
