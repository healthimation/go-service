package eventbus

import (
	"context"
)

type Subscriber interface {
	Subscribe(key string, fn MessageHandler)
	Start(ctx context.Context) error
}

type SubscriberConfig struct {
	Source         string
	EventBusName   string
	QueueUrl       string
	MaxWorker      int
	MaxMsg         int
	DefaultHandler MessageHandler
}
