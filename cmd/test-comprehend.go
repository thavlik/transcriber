package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pkg/comprehend"
	"go.uber.org/zap"
)

var testComprehendCmd = &cobra.Command{
	Use:   "comprehend",
	Short: "test Amazon Comprehend with a text string",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		text, ok := os.LookupEnv("TEXT")
		if !ok {
			if len(args) == 0 {
				return errors.New("no text provided")
			}
			text = strings.Join(args, " ")
		}
		DefaultLog.Info("testing comprehend", zap.String("text", text))
		entities, err := comprehend.Comprehend(
			context.Background(),
			text,
			DefaultLog,
		)
		if err != nil {
			return err
		}
		for _, entity := range entities {
			DefaultLog.Info(
				"entity",
				zap.String("text", entity.Text),
				zap.String("type", entity.Type),
				zap.Float64("score", entity.Score),
			)
		}
		return nil
	},
}

func init() {
	testCmd.AddCommand(testComprehendCmd)
}
