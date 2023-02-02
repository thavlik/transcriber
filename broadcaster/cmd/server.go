package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/broadcaster/pkg/server"
)

var serverArgs struct {
	httpPort    int
	metricsPort int
	redis       base.RedisOptions
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.RedisEnv(&serverArgs.redis, false)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var redisClient *redis.Client
		if serverArgs.redis.IsSet() {
			redisClient = base.ConnectRedis(
				cmd.Context(),
				&serverArgs.redis,
			)
		}
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.metricsPort,
			redisClient,
			base.DefaultLog,
		)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
}
