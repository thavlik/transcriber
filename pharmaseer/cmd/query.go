package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use: "query",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
