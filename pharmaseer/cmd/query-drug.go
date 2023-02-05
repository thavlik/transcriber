package main

import (
	"encoding/json"
	"fmt"

	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/pharmaseer/pkg/api"

	"github.com/spf13/cobra"
)

var queryDrugArgs struct {
	base.ServiceOptions
	force bool
}

var queryDrugCmd = &cobra.Command{
	Use:  "drug",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		seer := api.NewPharmaSeerClientFromOptions(queryDrugArgs.ServiceOptions)
		details, err := seer.GetDrugDetails(
			cmd.Context(),
			api.GetDrugDetails{
				Query: args[0],
				Force: queryDrugArgs.force,
			})
		if err != nil {
			return err
		}
		body, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	queryCmd.AddCommand(queryDrugCmd)
	base.AddServiceFlags(queryDrugCmd, "", &queryDrugArgs.ServiceOptions, 0)
	queryDrugCmd.PersistentFlags().BoolVarP(&queryDrugArgs.force, "force", "f", false, "force query from youtube")
}
