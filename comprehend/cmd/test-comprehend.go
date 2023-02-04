package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend"
	"github.com/thavlik/transcriber/comprehend/pkg/comprehend/adapter"
	"go.uber.org/zap"
)

var testComprehendArgs struct {
	service      string
	includeTerms string
	excludeTerms string
	includeTypes string
	excludeTypes string
	languageCode string
}

var testComprehendCmd = &cobra.Command{
	Use:   "comprehend",
	Short: "test Amazon Comprehend with a text string",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		text, ok := os.LookupEnv("TEXT")
		if !ok {
			text = strings.TrimSpace(strings.Join(args, " "))
			if len(text) == 0 {
				return errors.New("no text provided")
			}
		}
		base.DefaultLog.Info("testing comprehend", zap.String("text", text))
		entities, err := adapter.Comprehend(
			cmd.Context(),
			adapter.Model(testComprehendArgs.service),
			text,
			testComprehendArgs.languageCode,
			&comprehend.Filter{
				IncludeTerms: strings.Split(testComprehendArgs.includeTerms, ","),
				ExcludeTerms: strings.Split(testComprehendArgs.excludeTerms, ","),
				IncludeTypes: strings.Split(testComprehendArgs.includeTypes, ","),
				ExcludeTypes: strings.Split(testComprehendArgs.excludeTypes, ","),
			},
		)
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(entities, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	testCmd.AddCommand(testComprehendCmd)
	testComprehendCmd.Flags().StringVarP(&testComprehendArgs.service, "service", "s", "default", "Amazon Comprehend service (default, medical)")
	testComprehendCmd.Flags().StringVar(&testComprehendArgs.includeTypes, "include-types", "", "entity type include filter (comma-separated list)")
	testComprehendCmd.Flags().StringVar(&testComprehendArgs.excludeTypes, "exclude-types", "", "entity type exclude filter (comma-separated list)")
	testComprehendCmd.Flags().StringVar(&testComprehendArgs.includeTerms, "include-terms", "", "entity type include filter (comma-separated list, case sensitive)")
	testComprehendCmd.Flags().StringVar(&testComprehendArgs.excludeTerms, "exclude-terms", "", "entity type exclude filter (comma-separated list, case sensitive)")
	testComprehendCmd.Flags().StringVarP(&testComprehendArgs.languageCode, "language-code", "l", "en", "shorthand language code (default is 'en' for English)")
}
