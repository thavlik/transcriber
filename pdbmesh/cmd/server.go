package main

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	memory_scheduler "github.com/thavlik/transcriber/base/pkg/scheduler/memory"
	redis_scheduler "github.com/thavlik/transcriber/base/pkg/scheduler/redis"
	"github.com/thavlik/transcriber/pdbmesh/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	redis    base.RedisOptions
	s3Bucket string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.CheckEnv("S3_BUCKET", &serverArgs.s3Bucket)
		if serverArgs.s3Bucket == "" {
			return errors.New("missing --s3-bucket")
		}
		base.RedisEnv(&serverArgs.redis, false)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		redis := initRedis(cmd.Context())
		return server.Entry(
			cmd.Context(),
			&serverArgs.ServerOptions,
			initScheduler(redis, "pdbmesh"),
			base.DefaultLog,
		)
	},
}

func initRedis(ctx context.Context) *redis.Client {
	if serverArgs.redis.IsSet() {
		return base.ConnectRedis(ctx, &serverArgs.redis)
	}
	return nil
}

func initScheduler(
	redis *redis.Client,
	name string,
) scheduler.Scheduler {
	if redis != nil {
		return redis_scheduler.NewRedisScheduler(
			redis,
			name,
			10*time.Second,
		)
	}
	return memory_scheduler.NewMemoryScheduler()
}

func init() {
	rootCmd.AddCommand(serverCmd)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.Flags().StringVar(&serverArgs.s3Bucket, "s3-bucket", "", "name of the s3 bucket to store image data in")
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
}
