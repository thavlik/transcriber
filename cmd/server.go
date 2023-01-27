package main

import (
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/pkg/server"
)

var serverArgs struct {
	httpPort    int
	rtmpPort    int
	metricsPort int
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			serverArgs.httpPort,
			serverArgs.rtmpPort,
			serverArgs.metricsPort,
			DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.rtmpPort, "rtmp-port", "r", 1935, "rtmp port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
}
