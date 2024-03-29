package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	serverOpts *base.ServerOptions,
	iam iam.IAM,
	pubSub pubsub.PubSub,
	corsHeader string,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		iam,
		pubsub.Publisher(pubSub),
		corsHeader,
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			ctx,
			serverOpts.MetricsPort,
			log,
		)
	})

	sub, err := pubSub.Subscribe(
		s.ctx,
		channelName,
	)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe")
	}
	msgs := sub.Messages(s.ctx)
	s.spawn(func() {
		s.pumpMessages(msgs)
	})

	base.SignalReady(log)

	return s.ListenAndServe(
		serverOpts.Port,
	)
}
