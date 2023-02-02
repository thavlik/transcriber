package memory_pubsub

import "context"

func (p *memoryPubSub) Publish(
	ctx context.Context,
	topic string,
	payload []byte,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.l <- struct{}{}:
		defer func() { <-p.l }()
	}
	t, ok := p.channels[topic]
	if !ok {
		// no subscribers to this topic
		return nil
	}
	done := ctx.Done()
	for sub := range t {
		select {
		case <-done:
			return ctx.Err()
		case sub.ch <- payload:
		default:
			p.log.Warn("memory pubsub dropped message due to channel being full")
		}
	}
	return nil
}
