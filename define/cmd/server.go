package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/server"
)

var serverArgs struct {
	httpPort        int
	metricsPort     int
	openAISecretKey string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("OPENAI_SECRET_KEY", &serverArgs.openAISecretKey)
		if serverArgs.openAISecretKey == "" {
			return errors.New("missing --openai-secret-key")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.metricsPort,
			serverArgs.openAISecretKey,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVar(&serverArgs.openAISecretKey, "openai-secret-key", "", "OpenAI API secret key")
}
