package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/disease"
)

var testIsDiseaseArgs struct {
	openAISecretKey string
}

var testIsDiseaseCmd = &cobra.Command{
	Use:   "is-disease",
	Short: "test querying OpenAI to determine if an input string is a disease",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("OPENAI_SECRET_KEY", &testIsDiseaseArgs.openAISecretKey)
		if testIsDiseaseArgs.openAISecretKey == "" {
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
		client := gpt3.NewClient(
			testIsDiseaseArgs.openAISecretKey,
			gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
		)
		isDisease, err := disease.IsDisease(
			cmd.Context(),
			client,
			text,
		)
		if err != nil {
			return errors.Wrap(err, "disease.IsDisease")
		}
		fmt.Printf("%t", isDisease)
		return nil
	},
}

func init() {
	testCmd.AddCommand(testIsDiseaseCmd)
	testIsDiseaseCmd.Flags().StringVar(&serverArgs.openAISecretKey, "openai-secret-key", "", "OpenAI API secret key")
}
