package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pkg/source"
	"github.com/thavlik/transcriber/pkg/source/wav"
	"github.com/thavlik/transcriber/pkg/transcriber"
)

var testCmd = &cobra.Command{
	Use:  "test",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		var src source.Source
		switch filepath.Ext(input) {
		case ".wav":
			src, err = wav.NewWavSource(f)
			if err != nil {
				return errors.Wrap(err, "wav.NewWavSource")
			}
		default:
			return errors.New("unsupported file type")
		}
		return transcriber.Transcribe(
			context.Background(),
			src,
		)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
