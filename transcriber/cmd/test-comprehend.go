package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/transcriber/pkg/comprehend"
	"go.uber.org/zap"
)

var testComprehendArgs struct {
	service string
	include string
	exclude string
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
		var entities []*comprehend.Entity
		switch testComprehendArgs.service {
		case "default":
			entities, err = comprehend.Comprehend(
				cmd.Context(),
				text,
				strings.Split(testComprehendArgs.include, ","),
				strings.Split(testComprehendArgs.exclude, ","),
				base.DefaultLog,
			)
		case "medical":
			entities, err = comprehend.ComprehendMedical(
				cmd.Context(),
				text,
				strings.Split(testComprehendArgs.include, ","),
				strings.Split(testComprehendArgs.exclude, ","),
				base.DefaultLog,
			)
		default:
			return fmt.Errorf("invalid service '%s'", testComprehendArgs.service)
		}
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
	testComprehendCmd.Flags().StringVarP(&testComprehendArgs.include, "include", "i", "", "entity type include filter (comma-separated list)")
	testComprehendCmd.Flags().StringVarP(&testComprehendArgs.exclude, "exclude", "e", "", "entity type exclude filter (comma-separated list)")
}
