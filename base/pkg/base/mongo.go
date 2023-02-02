package base

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions struct {
	DBName     string `json:"dbName"`
	Host       string `json:"host"`
	AuthSource string `json:"authSource"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

func (o *MongoOptions) IsSet() bool {
	return o.Host != ""
}

func (o *MongoOptions) Ensure() {
	if o.DBName == "" {
		panic("missing --mongo-db-name")
	}
	if o.Host == "" {
		panic("missing --mongo-host")
	}
	if o.AuthSource == "" {
		panic("missing --mongo-auto-source")
	}
	if o.Username == "" {
		panic("missing --mongo-username")
	}
	if o.Password == "" {
		panic("missing --mongo-password")
	}
}

func AddMongoFlags(cmd *cobra.Command, o *MongoOptions) {
	cmd.PersistentFlags().StringVar(&o.DBName, "mongo-db-name", "ts", "mongodb database name")
	cmd.PersistentFlags().StringVar(&o.Host, "mongo-host", "", "mongodb service endpoint, with port")
	cmd.PersistentFlags().StringVar(&o.AuthSource, "mongo-auth-source", "admin", "mongodb auth db name")
	cmd.PersistentFlags().StringVar(&o.Username, "mongo-username", "", "mongodb username")
	cmd.PersistentFlags().StringVar(&o.Password, "mongo-password", "", "mongodb password")
}

func ConnectMongo(
	ctx context.Context,
	o *MongoOptions,
) *mongo.Database {
	if !o.IsSet() {
		panic("missing mongo options")
	}
	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI("mongodb+srv://"+o.Host).
			SetAuth(options.Credential{
				AuthSource: o.AuthSource,
				Username:   o.Username,
				Password:   o.Password,
			}))
	if err != nil {
		panic(fmt.Errorf("mongo: error connecting to %s: %v", o.Host, err))
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("mongo: failed to ping: %v", err))
	}
	db := mongoClient.Database(o.DBName)
	DefaultLog.Debug("connected to mongo", Elapsed(start))
	return db
}

func MongoEnv(o *MongoOptions, required bool) *MongoOptions {
	if o == nil {
		o = &MongoOptions{}
	}
	CheckEnv("MONGO_DB_NAME", &o.DBName)
	CheckEnv("MONGO_HOST", &o.Host)
	CheckEnv("MONGO_AUTH_SOURCE", &o.AuthSource)
	CheckEnv("MONGO_USERNAME", &o.Username)
	CheckEnv("MONGO_PASSWORD", &o.Password)
	if required {
		o.Ensure()
	}
	return o
}
