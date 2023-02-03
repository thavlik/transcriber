package server

import (
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func Entry(
	port int,
	adminPort int,
	iam iam.IAM,
	imgSearch *base.ServiceOptions,
	corsHeader string,
	log *zap.Logger,
) error {
	s := NewServer(
		iam,
		imgSearch,
		corsHeader,
		log,
	)
	mainErr := make(chan error, 1)
	go func() { mainErr <- s.ListenAndServe(port) }()
	adminErr := make(chan error, 1)
	go func() { adminErr <- s.AdminListenAndServe(adminPort) }()
	select {
	case err := <-mainErr:
		return err
	case err := <-adminErr:
		return err
	}
}