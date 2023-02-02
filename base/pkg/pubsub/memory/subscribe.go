package memory_pubsub

import (
	"context"

	"github.com/thavlik/transcriber/base/pkg/pubsub"
)

func (p *memoryPubSub) Subscribe(
	ctx context.Context,
	topic string,
) (pubsub.Subscription, error) {
	sub := &memorySubscription{
		topic: topic,
		ch:    make(chan []byte, 32),
		p:     p,
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case p.l <- struct{}{}:
		defer func() { <-p.l }()
	}
	t, ok := p.channels[topic]
	if !ok {
		t = make(map[*memorySubscription]struct{})
		p.channels[topic] = t
	}
	t[sub] = struct{}{}
	return sub, nil
}
