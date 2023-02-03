package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "define",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}
