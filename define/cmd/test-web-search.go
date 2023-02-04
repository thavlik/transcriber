package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/websearch/adapter"
	"go.uber.org/zap"
)

var testWebSearchArgs struct {
	input    string
	service  string
	endpoint string
	apiKey   string
	count    int
	offset   int
}

var testWebSearchCmd = &cobra.Command{
	Use:   "web-search",
	Short: "test web search with a text string",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("BING_API_KEY", &testWebSearchArgs.apiKey)
		if testWebSearchArgs.apiKey == "" {
			return errors.New("BING_API_KEY not set")
		}
		base.CheckEnv("BING_ENDPOINT", &testWebSearchArgs.endpoint)
		base.CheckEnv("INPUT", &testWebSearchArgs.input)
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
		service := adapter.Bing
		base.DefaultLog.Info(
			"testing web search",
			zap.String("service", string(service)),
			zap.String("text", text))
		results, err := adapter.Search(
			cmd.Context(),
			service,
			text,
			testWebSearchArgs.endpoint,
			testWebSearchArgs.apiKey,
			testWebSearchArgs.count,
			testWebSearchArgs.offset,
		)
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	testCmd.AddCommand(testWebSearchCmd)
	testWebSearchCmd.Flags().StringVarP(
		&testWebSearchArgs.service,
		"service",
		"s",
		"bing",
		"service name (bing only for now)",
	)
	testWebSearchCmd.Flags().StringVarP(
		&testWebSearchArgs.input,
		"input",
		"i",
		"",
		"input text",
	)
	testWebSearchCmd.Flags().StringVarP(
		&testWebSearchArgs.endpoint,
		"bing-endpoint",
		"e",
		defaultBingEndpoint,
		"bing search endpoint",
	)
	testWebSearchCmd.Flags().StringVarP(
		&testWebSearchArgs.apiKey,
		"bing-api-key",
		"k",
		"",
		"bing api key",
	)
	testWebSearchCmd.Flags().IntVarP(
		&testWebSearchArgs.count,
		"count",
		"c",
		10,
		"count",
	)
	testWebSearchCmd.Flags().IntVarP(
		&testWebSearchArgs.offset,
		"offset",
		"o",
		0,
		"offset",
	)
}
