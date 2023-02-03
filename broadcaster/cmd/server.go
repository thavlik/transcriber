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
	base.ServerOptions
	iam        base.IAMOptions
	redis      base.RedisOptions
	corsHeader string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.RedisEnv(&serverArgs.redis, false)
		base.IAMEnv(&serverArgs.iam, false)
		base.CheckEnv("CORS_HEADER", &serverArgs.corsHeader)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			&serverArgs.ServerOptions,
			nil, //iam.InitIAM(&serverArgs.iam, base.DefaultLog),
			initPubSub(
				cmd.Context(),
				base.DefaultLog,
			),
			serverArgs.corsHeader,
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
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.Flags().StringVar(&serverArgs.corsHeader, "cors-header", "", "Access-Control-Allow-Origin header")
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
	base.AddIAMFlags(serverCmd, &serverArgs.iam)
}
