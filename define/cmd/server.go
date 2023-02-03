package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/define/pkg/server"
	"github.com/thavlik/transcriber/define/pkg/storage"
	mongo_storage "github.com/thavlik/transcriber/define/pkg/storage/mongo"
)

var serverArgs struct {
	base.ServerOptions
	openAISecretKey string
	db              base.DatabaseOptions
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServerEnv(&serverArgs.ServerOptions)
		base.DatabaseEnv(&serverArgs.db, true)
		base.CheckEnv("OPENAI_SECRET_KEY", &serverArgs.openAISecretKey)
		if serverArgs.openAISecretKey == "" {
			return errors.New("missing --openai-secret-key")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Entry(
			cmd.Context(),
			&serverArgs.ServerOptions,
			initStorage(cmd.Context()),
			serverArgs.openAISecretKey,
			base.DefaultLog,
		)
	},
}

func initStorage(ctx context.Context) storage.Storage {
	switch serverArgs.db.Driver {
	case "mongo":
		return mongo_storage.NewMongoStorage(
			base.ConnectMongo(
				ctx,
				&serverArgs.db.Mongo,
			).Collection("definitions"))
	default:
		panic(fmt.Errorf("unsupported storage driver '%s'", serverArgs.db.Driver))
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	serverCmd.Flags().StringVar(&serverArgs.openAISecretKey, "openai-secret-key", "", "OpenAI API secret key")
	base.AddDatabaseFlags(serverCmd, &serverArgs.db)
}
