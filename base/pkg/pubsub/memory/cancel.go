package memory_pubsub

import "context"

func (m *memoryPubSub) cancel(
	ctx context.Context,
	sub *memorySubscription,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.l <- struct{}{}:
		defer func() { <-m.l }()
	}
	t, ok := m.channels[sub.topic]
	if !ok {
		panic("memory pubsub: canceling subscription for unknown topic")
	}
	delete(t, sub)
	if len(t) == 0 {
		// no more subscriptions for this topic
		delete(m.channels, sub.topic)
	}
	close(sub.ch)
	return nil
}
