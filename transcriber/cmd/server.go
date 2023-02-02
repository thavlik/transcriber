package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/transcriber/pkg/server"
)

var serverArgs struct {
	httpPort    int
	metricsPort int
	streamKey   string
	specialty   string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if v, ok := os.LookupEnv("STREAM_KEY"); ok {
			serverArgs.streamKey = v
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.metricsPort,
			serverArgs.specialty,
			serverArgs.streamKey,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVarP(&serverArgs.streamKey, "stream-key", "k", "", "stream key to use for authentication")
	serverCmd.Flags().StringVarP(
		&serverArgs.specialty,
		"specialty",
		"s",
		defaultSpecialty,
		"the specialty to use for transcription",
	)
}
