package disease

import (
	"context"
	"fmt"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
)

var (
	temp float32 = 0.7
	topP float32 = 1.0
)

func IsDisease(
	ctx context.Context,
	client gpt3.Client,
	input string,
) (bool, error) {
	input = isDiseaseQuery(input)
	n := 1
	maxLength := 5 // small number for yes or no answers
	resp, err := client.Completion(
		ctx,
		gpt3.CompletionRequest{
			Prompt:           []string{input},
			Temperature:      &temp,
			MaxTokens:        &maxLength,
			TopP:             &topP,
			N:                &n,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
		},
	)
	if err != nil {
		return false, errors.Wrap(err, "gpt3")
	}
	for _, choice := range resp.Choices {
		output := strings.ToLower(strings.TrimSpace(choice.Text))
		if strings.Contains(output, "yes") {
			return true, nil
		} else if strings.Contains(output, "no") {
			return false, nil
		} else {
			continue
		}
	}
	return false, errors.New("invalid response from gpt3")
}

func isDiseaseQuery(input string) string {
	return fmt.Sprintf(
		"Yes or no, is the term \"%s\" a disease?",
		input,
	)
}
