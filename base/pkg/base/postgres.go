package base

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

type PostgresOptions struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	SSLMode    string
	SSLKeyPath string
}

func (p *PostgresOptions) IsSet() bool {
	return p.Host != ""
}

func (o *PostgresOptions) Ensure() {
	if o.Host == "" {
		panic("missing --postgres-host")
	}
	if o.Username == "" {
		panic("missing --postgres-username")
	}
	if o.Password == "" {
		panic("missing --postgres-password")
	}
	if o.Database == "" {
		panic("missing --postgres-database")
	}
}

func AddPostgresFlags(cmd *cobra.Command, out *PostgresOptions) {
	cmd.PersistentFlags().StringVar(&out.Host, "postgres-host", "", "PostgreSQL hostname")
	cmd.PersistentFlags().IntVar(&out.Port, "postgres-port", 5432, "PostgreSQL port")
	cmd.PersistentFlags().StringVar(&out.Username, "postgres-username", "", "PostgreSQL username")
	cmd.PersistentFlags().StringVar(&out.Password, "postgres-password", "", "PostgreSQL password")
	cmd.PersistentFlags().StringVar(&out.Database, "postgres-database", "", "PostgreSQL database name")
	cmd.PersistentFlags().StringVar(&out.SSLMode, "postgres-ssl-mode", "disable", "PostgreSQL SSL mode")
	cmd.PersistentFlags().StringVar(&out.SSLKeyPath, "postgres-ssl-key-path", "", "PostgreSQL CA cert path")
}

func PostgresEnv(o *PostgresOptions, required bool) *PostgresOptions {
	if o == nil {
		o = &PostgresOptions{}
	}

	CheckEnv("POSTGRES_HOST", &o.Host)
	CheckEnvInt("POSTGRES_PORT", &o.Port)
	CheckEnv("POSTGRES_USERNAME", &o.Username)
	CheckEnv("POSTGRES_PASSWORD", &o.Password)
	CheckEnv("POSTGRES_DATABASE", &o.Database)
	CheckEnv("POSTGRES_SSL_MODE", &o.SSLMode)
	CheckEnv("POSTGRES_SSL_KEY_PATH", &o.SSLKeyPath)

	if required {
		o.Ensure()
	}

	// Allow cert injection from env, overriding SSLKeyPath
	caCert := os.Getenv("POSTGRES_CA_CERT")
	if len(caCert) > 0 {
		o.SSLKeyPath = "/etc/cacert.pem"
		if err := os.WriteFile(
			o.SSLKeyPath,
			[]byte(caCert),
			0644,
		); err != nil {
			panic(fmt.Errorf("write %s: %v", o.SSLKeyPath, err))
		}
	}

	return o
}

func ConnectPostgres(o *PostgresOptions) *sql.DB {
	if !o.IsSet() {
		panic("missing postgres options")
	}
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		o.Host,
		o.Port,
		o.Username,
		o.Password,
		o.Database,
		o.SSLMode,
	)
	if o.SSLKeyPath != "" {
		psqlInfo += fmt.Sprintf(" sslkey=%s", o.SSLKeyPath)
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(fmt.Errorf("sql.Open: %v", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("ping postgres: %v", err))
	}
	DefaultLog.Debug("connected to postgres", Elapsed(start))
	return db
}
