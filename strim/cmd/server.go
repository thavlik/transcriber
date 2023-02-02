package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/strim/pkg/server"
)

var serverArgs struct {
	httpPort    int
	rtmpPort    int
	metricsPort int
	transcriber base.ServiceOptions
	streamKey   string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServiceEnv("transcriber", &serverArgs.transcriber)
		if v, ok := os.LookupEnv("STREAM_KEY"); ok {
			serverArgs.streamKey = v
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.rtmpPort,
			serverArgs.metricsPort,
			serverArgs.transcriber,
			serverArgs.streamKey,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	base.AddServiceFlags(serverCmd, "transcriber", &serverArgs.transcriber, 6*time.Second)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.rtmpPort, "rtmp-port", "r", 1935, "rtmp port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVarP(&serverArgs.streamKey, "stream-key", "k", "", "stream key to use for authentication")
}
