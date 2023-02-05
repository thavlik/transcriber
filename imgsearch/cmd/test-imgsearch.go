package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/imgsearch/adapter"
)

var testImgSearchArgs struct {
	query    string
	service  string
	endpoint string
	apiKey   string
	count    int
	offset   int
}

var testImgSearch = &cobra.Command{
	Use:  "imgsearch",
	Args: cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("BING_API_KEY", &testImgSearchArgs.apiKey)
		if testImgSearchArgs.apiKey == "" {
			return errors.New("BING_API_KEY not set")
		}
		base.CheckEnv("BING_ENDPOINT", &testImgSearchArgs.endpoint)
		base.CheckEnv("QUERY", &testImgSearchArgs.query)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if testImgSearchArgs.query == "" {
			testImgSearchArgs.query = strings.TrimSpace(strings.Join(args, " "))
			if len(testImgSearchArgs.query) == 0 {
				return errors.New("no text provided")
			}
		}
		result, err := adapter.Search(
			cmd.Context(),
			adapter.SearchService(testImgSearchArgs.service),
			testImgSearchArgs.query,
			testImgSearchArgs.endpoint,
			testImgSearchArgs.apiKey,
			testImgSearchArgs.count,
			testImgSearchArgs.offset,
		)
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	testCmd.AddCommand(testImgSearch)
	testImgSearch.Flags().StringVarP(
		&testImgSearchArgs.service,
		"service",
		"s",
		"bing",
		"service name (bing only for now)",
	)
	testImgSearch.Flags().StringVarP(
		&testImgSearchArgs.query,
		"query",
		"q",
		"",
		"input text query",
	)
	testImgSearch.Flags().StringVarP(
		&testImgSearchArgs.endpoint,
		"bing-endpoint",
		"e",
		defaultBingEndpoint,
		"bing search endpoint",
	)
	testImgSearch.Flags().StringVarP(
		&testImgSearchArgs.apiKey,
		"bing-api-key",
		"k",
		"",
		"bing api key",
	)
	testImgSearch.Flags().IntVarP(
		&testImgSearchArgs.count,
		"count",
		"c",
		10,
		"count",
	)
	testImgSearch.Flags().IntVarP(
		&testImgSearchArgs.offset,
		"offset",
		"o",
		0,
		"offset",
	)
}
