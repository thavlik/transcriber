package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/data/s3"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/meta"
	"github.com/thavlik/transcriber/imgsearch/pkg/cache/meta/mongo"
	"github.com/thavlik/transcriber/imgsearch/pkg/server"
)

var serverArgs struct {
	httpPort           int
	metricsPort        int
	bingApiKey         string
	bingEndpoint       string
	db                 base.DatabaseOptions
	metaCollectionName string
	s3Bucket           string
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
		return server.Entry(
			serverArgs.httpPort,
			serverArgs.metricsPort,
			serverArgs.bingApiKey,
			serverArgs.bingEndpoint,
			cache.NewImageCache(
				initMetaCache(),
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

func initMetaCache() meta.ImageMetaCache {
	return mongo.NewMongoMetaCache(
		base.ConnectMongo(
			&serverArgs.db.Mongo,
		).Collection(serverArgs.metaCollectionName),
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
	serverCmd.Flags().StringVar(&serverArgs.s3Bucket, "s3-bucket", "", "name of the s3 bucket to store image data in")
}
