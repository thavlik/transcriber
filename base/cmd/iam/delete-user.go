package iam

import (
	"context"
	"errors"
	"time"

	"github.com/thavlik/transcriber/base/pkg/iam/api"

	"github.com/spf13/cobra"
)

var deleteUserArgs struct {
	endpoint       string
	timeout        time.Duration
	id             string
	username       string
	deleteProjects bool
}

var deleteUserCmd = &cobra.Command{
	Use: "delete-user",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if deleteUserArgs.id == "" && deleteUserArgs.username == "" {
			return errors.New("missing --id or --username")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := api.NewRemoteIAMClient(
			deleteUserArgs.endpoint,
			api.NewRemoteIAMClientOptions().
				SetTimeout(deleteUserArgs.timeout),
		).DeleteUser(
			context.Background(),
			api.DeleteUser{
				ID:             deleteUserArgs.id,
				Username:       deleteUserArgs.username,
				DeleteProjects: deleteUserArgs.deleteProjects,
			},
		)
		return err
	},
}

func init() {
	deleteUserCmd.PersistentFlags().StringVar(&deleteUserArgs.endpoint, "endpoint", "http://localhost:8080", "admin service endpoint")
	deleteUserCmd.PersistentFlags().DurationVar(&deleteUserArgs.timeout, "timeout", defaultTimeout, "admin service timeout")
	deleteUserCmd.PersistentFlags().StringVar(&deleteUserArgs.id, "id", "", "user id")
	deleteUserCmd.PersistentFlags().StringVar(&deleteUserArgs.username, "username", "", "username")
	deleteUserCmd.PersistentFlags().BoolVar(&deleteUserArgs.deleteProjects, "delete-projects", false, "also delete the user's projects")
	iamCmd.AddCommand(deleteUserCmd)
}
