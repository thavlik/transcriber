package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	memory_pubsub "github.com/thavlik/transcriber/base/pkg/pubsub/memory"
	redis_pubsub "github.com/thavlik/transcriber/base/pkg/pubsub/redis"
	"github.com/thavlik/transcriber/broadcaster/pkg/server"
	"go.uber.org/zap"
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
			initPubSub(
				cmd.Context(),
				base.DefaultLog,
			),
			base.DefaultLog,
		)
	},
}

func initPubSub(
	ctx context.Context,
	log *zap.Logger,
) pubsub.PubSub {
	if serverArgs.redis.IsSet() {
		return redis_pubsub.NewRedisPubSub(
			base.ConnectRedis(
				ctx,
				&serverArgs.redis,
			),
			log,
		)
	}
	return memory_pubsub.NewMemoryPubSub(log)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
}
