package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/strim/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	rtmpPort  int
	scribe    base.ServiceOptions
	streamKey string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.ServiceEnv("scribe", &serverArgs.scribe)
		if v, ok := os.LookupEnv("STREAM_KEY"); ok {
			serverArgs.streamKey = v
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			&serverArgs.ServerOptions,
			serverArgs.rtmpPort,
			serverArgs.scribe,
			serverArgs.streamKey,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	base.AddServiceFlags(serverCmd, "scribe", &serverArgs.scribe, 6*time.Second)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.Flags().IntVarP(&serverArgs.rtmpPort, "rtmp-port", "r", 1935, "rtmp port to listen on")
	serverCmd.Flags().StringVarP(&serverArgs.streamKey, "stream-key", "k", "", "stream key to use for authentication")
}
