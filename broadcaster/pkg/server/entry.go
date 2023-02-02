package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/pubsub"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	metricsPort int,
	pubSub pubsub.PubSub,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		pubsub.Publisher(pubSub),
		log,
	)
	defer s.ShutDown()

	s.spawn(func() {
		base.RunMetrics(
			ctx,
			metricsPort,
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

	return s.ListenAndServe(
		httpPort,
	)
}
