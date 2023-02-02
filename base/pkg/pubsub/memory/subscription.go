package memory_pubsub

import (
	"context"
)

type memorySubscription struct {
	topic string
	ch    chan []byte
	p     *memoryPubSub
}

func (m *memorySubscription) Messages(ctx context.Context) <-chan []byte {
	return m.ch
}

func (m *memorySubscription) Cancel(ctx context.Context) error {
	return m.p.cancel(ctx, m)
}
