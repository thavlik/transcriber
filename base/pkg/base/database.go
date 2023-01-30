package base

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type DatabaseDriver string

var (
	MongoDriver    DatabaseDriver = "mongo"
	PostgresDriver DatabaseDriver = "postgres"

	databaseDrivers []DatabaseDriver = []DatabaseDriver{MongoDriver, PostgresDriver}
)

func (d DatabaseDriver) IsValid() bool {
	for _, driver := range databaseDrivers {
		if driver == d {
			return true
		}
	}
	return false
}

type DatabaseOptions struct {
	Driver   DatabaseDriver
	Mongo    MongoOptions
	Postgres PostgresOptions
}

func (o *DatabaseOptions) IsSet() bool {
	switch o.Driver {
	case "":
		return false
	case MongoDriver:
		return o.Mongo.IsSet()
	case PostgresDriver:
		return o.Postgres.IsSet()
	default:
		panic(fmt.Errorf("unrecognized database driver '%s'", o.Driver))
	}
}

func AddDatabaseFlags(cmd *cobra.Command, o *DatabaseOptions) {
	AddMongoFlags(cmd, &o.Mongo)
	AddPostgresFlags(cmd, &o.Postgres)
	cmd.PersistentFlags().StringVar(((*string)(&o.Driver)), "db-driver", "", "database driver [ mongo | postgres ]")
}

func DatabaseEnv(o *DatabaseOptions, required bool) *DatabaseOptions {
	if o == nil {
		o = &DatabaseOptions{}
	}
	if v, ok := os.LookupEnv("DB_DRIVER"); ok {
		o.Driver = DatabaseDriver(v)
	}
	if !o.Driver.IsValid() {
		panic(fmt.Errorf("unrecognized database driver '%s'", o.Driver))
	}
	MongoEnv(&o.Mongo, required && o.Driver == MongoDriver)
	PostgresEnv(&o.Postgres, required && o.Driver == PostgresDriver)
	return o
}
