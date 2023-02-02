package server

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/thavlik/transcriber/base/pkg/base"
	"go.uber.org/zap"
)

func Entry(
	ctx context.Context,
	httpPort int,
	metricsPort int,
	redisClient *redis.Client,
	log *zap.Logger,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := NewServer(
		ctx,
		redisClient,
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

	if redisClient != nil {
		s.spawn(s.runSub)
	}

	return s.ListenAndServe(
		httpPort,
	)
}

func (s *Server) runSub() {
	ch := s.redisClient.Subscribe(
		s.ctx,
		channelName,
	).Channel()
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				panic("redis channel closed")
			}
			go s.broadcastLocal(
				s.ctx,
				[]byte(msg.Payload),
			)
		}
	}
}
