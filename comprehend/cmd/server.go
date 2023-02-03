package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/comprehend/pkg/entitycache"
	memory_entitycache "github.com/thavlik/transcriber/comprehend/pkg/entitycache/memory"
	redis_entitycache "github.com/thavlik/transcriber/comprehend/pkg/entitycache/redis"
	"github.com/thavlik/transcriber/comprehend/pkg/server"
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
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.metricsPort,
			initEntityCache(cmd.Context()),
			base.DefaultLog,
		)
	},
}

func initEntityCache(ctx context.Context) entitycache.EntityCache {
	if serverArgs.redis.IsSet() {
		return redis_entitycache.NewRedisEntityCache(
			base.ConnectRedis(ctx, &serverArgs.redis),
		)
	}
	return memory_entitycache.NewMemoryEntityCache()
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
}
