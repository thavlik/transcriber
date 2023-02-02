package pubsub

import "context"

type Subscription interface {
	Messages(ctx context.Context) <-chan []byte
	Cancel(ctx context.Context) error
}

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (Subscription, error)
}

type PubSub interface {
	Publisher
	Subscriber
}
