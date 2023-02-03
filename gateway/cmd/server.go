package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/gateway/pkg/server"
)

var serverArgs struct {
	base.ServerOptions
	adminPort  int
	iam        base.IAMOptions
	imgSearch  base.ServiceOptions
	define     base.ServiceOptions
	corsHeader string
}

var serverCmd = &cobra.Command{
	Use: "server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServiceEnv("imgsearch", &serverArgs.imgSearch)
		base.ServiceEnv("define", &serverArgs.define)
		base.ServerEnv(&serverArgs.ServerOptions)
		base.IAMEnv(&serverArgs.iam, false)
		base.CheckEnv("CORS_HEADER", &serverArgs.corsHeader)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log := base.DefaultLog
		return server.Entry(
			serverArgs.Port,
			serverArgs.adminPort,
			nil, //iam.InitIAM(&serverArgs.iam, log),
			&serverArgs.imgSearch,
			&serverArgs.define,
			serverArgs.corsHeader,
			log,
		)
	},
}

func init() {
	base.AddServerFlags(serverCmd, &serverArgs.ServerOptions)
	base.AddServiceFlags(serverCmd, "imgsearch", &serverArgs.imgSearch, 8*time.Second)
	base.AddServiceFlags(serverCmd, "define", &serverArgs.define, 8*time.Second)
	serverCmd.PersistentFlags().IntVar(&serverArgs.adminPort, "admin-port", 8080, "http service port")
	serverCmd.PersistentFlags().StringVar(&serverArgs.corsHeader, "cors-header", "", "Access-Control-Allow-Origin header")
	base.AddIAMFlags(serverCmd, &serverArgs.iam)
	ConfigureCommand(serverCmd)
}
