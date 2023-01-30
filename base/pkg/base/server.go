package base

import (
	"github.com/spf13/cobra"
)

type ServerOptions struct {
	Port        int
	MetricsPort int
}

func AddServerFlags(cmd *cobra.Command, out *ServerOptions) {
	cmd.PersistentFlags().IntVarP(&out.Port, "port", "p", 80, "http service port")
	cmd.PersistentFlags().IntVar(&out.MetricsPort, "metrics-port", 0, "prometheus metrics server scrape port")
}

func ServerEnv(o *ServerOptions) *ServerOptions {
	if o == nil {
		o = &ServerOptions{}
	}
	CheckEnvInt("PORT", &o.Port)
	CheckEnvInt("METRICS_PORT", &o.MetricsPort)
	return o
}
