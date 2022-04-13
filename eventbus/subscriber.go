package eventbus

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type Subscriber interface {
	Subscribe(ctx context.Context, key string, fn MessageHandler) error
	Start(ctx context.Context) error
}

type eventBridgeSubscriber struct {
	clt          *eventbridge.EventBridge
	source       string
	eventBusName string
}

func NewEventBridgeSubscriber(sess *session.Session, cfg Config) (Subscriber, error) {
	svc := eventbridge.New(sess)

	return &eventBridgeSubscriber{
		clt:          svc,
		source:       cfg.Source,
		eventBusName: cfg.EventBusName,
	}, nil
}

func (pub *eventBridgeSubscriber) Subscribe(ctx context.Context, key string, fn MessageHandler) error {
	return nil
}

func (pub *eventBridgeSubscriber) Start(ctx context.Context) error {
	return nil
}
