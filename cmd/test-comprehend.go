package main

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pkg/comprehend"
	"go.uber.org/zap"
)

var testComprehendCmd = &cobra.Command{
	Use:  "comprehend",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		text, ok := os.LookupEnv("TEXT")
		if !ok {
			return errors.New("TEXT environment variable not set")
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
