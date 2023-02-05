package main

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	memory_pubsub "github.com/thavlik/transcriber/base/pkg/pubsub/memory"
	redis_pubsub "github.com/thavlik/transcriber/base/pkg/pubsub/redis"
	"github.com/thavlik/transcriber/base/pkg/scheduler"
	memory_scheduler "github.com/thavlik/transcriber/base/pkg/scheduler/memory"
	redis_scheduler "github.com/thavlik/transcriber/base/pkg/scheduler/redis"
	"github.com/thavlik/transcriber/pharmaseer/pkg/infocache"
	mongo_infocache "github.com/thavlik/transcriber/pharmaseer/pkg/infocache/mongo"
	"github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache"
	s3_pdbcache "github.com/thavlik/transcriber/pharmaseer/pkg/pdbcache/s3"
	"github.com/thavlik/transcriber/pharmaseer/pkg/server"
	"go.uber.org/zap"
)

var serverArgs struct {
	base.ServerOptions
	redis       base.RedisOptions
	db          base.DatabaseOptions
	pdbBucket   string
	concurrency int
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.DatabaseEnv(&serverArgs.db, true)
		base.RedisEnv(&serverArgs.redis, false)
		base.CheckEnv("PDB_BUCKET", &serverArgs.pdbBucket)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		redis := initRedis(cmd.Context())
		return server.Entry(
			cmd.Context(),
			&serverArgs.ServerOptions,
			initPubSub(redis, log),
			initScheduler(redis, "dbqsched"),
			initScheduler(redis, "pdbsched"),
			initInfoCache(cmd.Context(), &serverArgs.db),
			initPDBCache(log),
			serverArgs.concurrency,
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

func initPubSub(
	redis *redis.Client,
	log *zap.Logger,
) pubsub.PubSub {
	if redis != nil {
		return redis_pubsub.NewRedisPubSub(
			redis,
			log,
		)
	}
	return memory_pubsub.NewMemoryPubSub(log)
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

func initPDBCache(log *zap.Logger) pdbcache.PDBCache {
	if serverArgs.pdbBucket != "" {
		return s3_pdbcache.NewS3PDBCache(serverArgs.pdbBucket, log)
	} else {
		panic(errors.New("missing pdb cache source"))
	}
}

func initInfoCache(
	ctx context.Context,
	opts *base.DatabaseOptions,
) infocache.InfoCache {
	switch opts.Driver {
	case "":
		panic("missing --db-driver")
	case base.MongoDriver:
		return mongo_infocache.NewMongoInfoCache(
			base.ConnectMongo(ctx, &opts.Mongo))
	default:
		panic(base.Unreachable)
	}
}

func init() {
	base.AddRedisFlags(serverCmd, &serverArgs.redis)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	serverCmd.PersistentFlags().IntVar(&serverArgs.concurrency, "concurrency", 1, "number of concurrent queries (best set to 1 and increase # replicas)")
	serverCmd.PersistentFlags().StringVar(&serverArgs.pdbBucket, "pdb-bucket", "", "Protein Data Bank (pdb) file cache bucket name")
	ConfigureCommand(serverCmd)
}
