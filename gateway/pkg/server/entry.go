package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	adminPort int,
	iam iam.IAM,
	imgSearch *base.ServiceOptions,
	define *base.ServiceOptions,
	corsHeader string,
	log *zap.Logger,
) error {
	s := NewServer(
		ctx,
		iam,
		imgSearch,
		define,
		corsHeader,
		log,
	)
	defer s.ShutDown()
	s.spawn(func() {
		base.RunMetrics(
			s.ctx,
			serverOpts.MetricsPort,
			log,
		)
	})
	mainErr := make(chan error, 1)
	s.spawn(func() { mainErr <- s.ListenAndServe(serverOpts.Port) })
	adminErr := make(chan error, 1)
	s.spawn(func() { adminErr <- s.AdminListenAndServe(adminPort) })
	base.SignalReady(log)
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case err := <-mainErr:
		return errors.Wrap(err, "main server error")
	case err := <-adminErr:
		return errors.Wrap(err, "admin server error")
	}
}
