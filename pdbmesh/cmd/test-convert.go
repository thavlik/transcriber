package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pdbmesh/pkg/convert"
)

var testConvertArgs struct {
	inputPath string
}

var testConvert = &cobra.Command{
	Use:  "convert",
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if testConvertArgs.inputPath == "" {
			testConvertArgs.inputPath = strings.TrimSpace(strings.Join(args, " "))
			if len(testConvertArgs.inputPath) == 0 {
				return errors.New("no text provided")
			}
		}
		f, err := os.Open(testConvertArgs.inputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		model, err := convert.Convert(
			cmd.Context(),
			f,
		)
		if err != nil {
			return err
		}
		defer model.Dispose()
		out, err := os.OpenFile("out.stl", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := io.Copy(out, model.Reader()); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	testCmd.AddCommand(testConvert)
	testConvert.Flags().StringVarP(
		&testConvertArgs.inputPath,
		"input",
		"i",
		"",
		"input pdb file path",
	)
}
