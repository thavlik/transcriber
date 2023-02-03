package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

var testCompletionArgs struct {
	openAISecretKey string
}

var testCompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "test OpenAI completion with a text string",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("OPENAI_SECRET_KEY", &testCompletionArgs.openAISecretKey)
		if testCompletionArgs.openAISecretKey == "" {
			return errors.New("missing --openai-secret-key")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		text, ok := os.LookupEnv("TEXT")
		if !ok {
			text = strings.TrimSpace(strings.Join(args, " "))
			if len(text) == 0 {
				return errors.New("no text provided")
			}
		}
		base.DefaultLog.Info("testing completion", zap.String("text", text))
		client := gpt3.NewClient(
			testCompletionArgs.openAISecretKey,
			gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
		)
		n := 1
		var temp float32 = 0.7
		var topP float32 = 1.0
		maxLength := 256
		resp, err := client.Completion(
			cmd.Context(),
			gpt3.CompletionRequest{
				Prompt:           []string{text},
				Temperature:      &temp,
				MaxTokens:        &maxLength,
				TopP:             &topP,
				N:                &n,
				FrequencyPenalty: 0.0,
				PresencePenalty:  0.0,
			},
		)
		if err != nil {
			return errors.Wrap(err, "gpt3")
		}
		body, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return errors.Wrap(err, "json")
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	testCmd.AddCommand(testCompletionCmd)
	testCompletionCmd.Flags().StringVar(&serverArgs.openAISecretKey, "openai-secret-key", "", "OpenAI API secret key")
}
