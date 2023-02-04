package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/search/adapter"
)

var testSearchArgs struct {
	input    string
	service  string
	endpoint string
	apiKey   string
	count    int
	offset   int
}

var testSearchCmd = &cobra.Command{
	Use:  "search",
	Args: cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("BING_API_KEY", &testSearchArgs.apiKey)
		if testSearchArgs.apiKey == "" {
			return errors.New("BING_API_KEY not set")
		}
		base.CheckEnv("BING_ENDPOINT", &testSearchArgs.endpoint)
		base.CheckEnv("INPUT", &testSearchArgs.input)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			testSearchArgs.input = strings.Join(args, " ")
		} else if testSearchArgs.input == "" {
			return errors.New("no search terms provided")
		}
		images, err := adapter.Search(
			cmd.Context(),
			adapter.SearchService(testSearchArgs.service),
			testSearchArgs.input,
			testSearchArgs.endpoint,
			testSearchArgs.apiKey,
			testSearchArgs.count,
			testSearchArgs.offset,
		)
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(images, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	testCmd.AddCommand(testSearchCmd)
	testSearchCmd.Flags().StringVarP(
		&testSearchArgs.service,
		"service",
		"s",
		"bing",
		"service name (bing only for now)",
	)
	testSearchCmd.Flags().StringVarP(
		&testSearchArgs.input,
		"input",
		"i",
		"",
		"input text",
	)
	testSearchCmd.Flags().StringVarP(
		&testSearchArgs.endpoint,
		"bing-endpoint",
		"e",
		defaultBingEndpoint,
		"bing search endpoint",
	)
	testSearchCmd.Flags().StringVarP(
		&testSearchArgs.apiKey,
		"bing-api-key",
		"k",
		"",
		"bing api key",
	)
	testSearchCmd.Flags().IntVarP(
		&testSearchArgs.count,
		"count",
		"c",
		10,
		"count",
	)
	testSearchCmd.Flags().IntVarP(
		&testSearchArgs.offset,
		"offset",
		"o",
		0,
		"offset",
	)
}
