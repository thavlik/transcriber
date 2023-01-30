package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/server"
)

var serverArgs struct {
	httpPort     int
	metricsPort  int
	bingApiKey   string
	bingEndpoint string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("BING_API_KEY", &serverArgs.bingApiKey)
		if serverArgs.bingApiKey == "" {
			return errors.New("BING_API_KEY not set")
		}
		base.CheckEnv("BING_ENDPOINT", &serverArgs.bingEndpoint)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			serverArgs.httpPort,
			serverArgs.metricsPort,
			serverArgs.bingApiKey,
			serverArgs.bingEndpoint,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVar(&serverArgs.bingApiKey, "bing-api-key", "", "bing api secret key")
	serverCmd.Flags().StringVar(&serverArgs.bingEndpoint, "bing-endpoint", defaultBingEndpoint, "bing search endpoint")
}
