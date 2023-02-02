package redis_pubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisSubscription struct {
	redis *redis.Client
	stop  chan struct{}
	topic string
	log   *zap.Logger
}

func (r *redisSubscription) Messages(ctx context.Context) <-chan []byte {
	sub := r.redis.Subscribe(ctx, r.topic).
		Channel(redis.WithChannelSize(64))
	ch := make(chan []byte, 32)
	go func() {
		defer close(ch)
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case <-r.stop:
				return
			case msg, ok := <-sub:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case <-r.stop:
					return
				case ch <- []byte(msg.Payload):
				default:
					r.log.Warn("redis subscription dropped message due to channel being full")
				}
			}
		}
	}()
	return ch
}

func (r *redisSubscription) Cancel(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case r.stop <- struct{}{}:
		return nil
	default:
		// already cancelled
		return nil
	}
}
