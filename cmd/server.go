package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pkg/server"
)

var serverArgs struct {
	httpPort    int
	rtmpPort    int
	metricsPort int
	streamKey   string
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
			serverArgs.httpPort,
			serverArgs.rtmpPort,
			serverArgs.metricsPort,
			serverArgs.streamKey,
			DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.rtmpPort, "rtmp-port", "r", 1935, "rtmp port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVarP(&serverArgs.streamKey, "stream-key", "s", "", "stream key to use for authentication")
}
