package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data/s3"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/meta"
	mongo_metacache "github.com/thavlik/transcriber/imgsearch/pkg/cache/meta/mongo"
	"github.com/thavlik/transcriber/imgsearch/pkg/history"
	mongo_history "github.com/thavlik/transcriber/imgsearch/pkg/history/mongo"
	"github.com/thavlik/transcriber/imgsearch/pkg/server"
	"go.mongodb.org/mongo-driver/mongo"
)

var serverArgs struct {
	httpPort              int
	metricsPort           int
	bingApiKey            string
	bingEndpoint          string
	db                    base.DatabaseOptions
	metaCollectionName    string
	historyCollectionName string
	s3Bucket              string
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.DatabaseEnv(&serverArgs.db, true)
		base.CheckEnv("META_COLLECTION_NAME", &serverArgs.metaCollectionName)
		if serverArgs.metaCollectionName == "" {
			return errors.New("missing --meta-collection-name")
		}
		base.CheckEnv("HISTORY_COLLECTION_NAME", &serverArgs.historyCollectionName)
		if serverArgs.historyCollectionName == "" {
			return errors.New("missing --history-collection-name")
		}
		base.CheckEnv("S3_BUCKET", &serverArgs.s3Bucket)
		if serverArgs.s3Bucket == "" {
			return errors.New("missing --s3-bucket")
		}
		base.CheckEnv("BING_API_KEY", &serverArgs.bingApiKey)
		if serverArgs.bingApiKey == "" {
			return errors.New("missing --bing-api-key")
		}
		base.CheckEnv("BING_ENDPOINT", &serverArgs.bingEndpoint)
		if serverArgs.bingEndpoint == "" {
			return errors.New("missing --bing-endpoint")
		}
		if serverArgs.db.Driver != "mongo" {
			return errors.New("only mongo is supported as a database driver")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mongo := base.ConnectMongo(cmd.Context(), &serverArgs.db.Mongo)
		return server.Entry(
			cmd.Context(),
			serverArgs.httpPort,
			serverArgs.metricsPort,
			initHistory(mongo),
			serverArgs.bingApiKey,
			serverArgs.bingEndpoint,
			cache.NewImageCache(
				initMetaCache(mongo),
				initDataCache(),
			),
			base.DefaultLog,
		)
	},
}

func initDataCache() data.ImageDataCache {
	return s3.NewS3DataCache(
		serverArgs.s3Bucket,
	)
}

func initMetaCache(db *mongo.Database) meta.ImageMetaCache {
	return mongo_metacache.NewMongoMetaCache(
		db.Collection(serverArgs.metaCollectionName),
	)
}

func initHistory(db *mongo.Database) history.History {
	return mongo_history.NewMongoHistory(
		db.Collection(serverArgs.historyCollectionName),
	)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverArgs.httpPort, "http-port", "p", 80, "http port to listen on")
	serverCmd.Flags().IntVarP(&serverArgs.metricsPort, "metrics-port", "m", 0, "metrics port to listen on")
	serverCmd.Flags().StringVar(&serverArgs.bingApiKey, "bing-api-key", "", "bing api secret key")
	serverCmd.Flags().StringVar(&serverArgs.bingEndpoint, "bing-endpoint", defaultBingEndpoint, "bing search endpoint")
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
	serverCmd.Flags().StringVar(&serverArgs.metaCollectionName, "meta-collection-name", "", "name of the collection to store image metadata in")
	serverCmd.Flags().StringVar(&serverArgs.historyCollectionName, "history-collection-name", "", "name of the collection to store image search history in")
	serverCmd.Flags().StringVar(&serverArgs.s3Bucket, "s3-bucket", "", "name of the s3 bucket to store image data in")
}
