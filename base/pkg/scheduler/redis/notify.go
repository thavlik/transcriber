package redis

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/base"
)

func (s *redisScheduler) Notify() <-chan struct{} {
	rc := s.redis.Subscribe(
		context.Background(),
		channelName(s.key),
	).Channel()
	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)
		for {
			msg, ok := <-rc
			if !ok {
				return
			} else if msg.Channel != channelName(s.key) {
				panic(base.Unreachable)
			}
			ch <- struct{}{}
		}
	}()
	return ch
}
