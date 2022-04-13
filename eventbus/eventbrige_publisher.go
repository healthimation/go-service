package eventbus

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"log"
	"time"
)

type PublisherConfig struct {
	Source       string
	EventBusName string
}

type eventBridgePublisher struct {
	clt          *eventbridge.EventBridge
	source       string
	eventBusName string
}

func NewEventBridgePublisher(sess *session.Session, cfg *PublisherConfig) (Publisher, error) {
	svc := eventbridge.New(sess)

	return &eventBridgePublisher{
		clt:          svc,
		source:       cfg.Source,
		eventBusName: cfg.EventBusName,
	}, nil
}

func (pub *eventBridgePublisher) Publish(ctx context.Context, msg *Message) error {
	resp, err := pub.clt.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				DetailType: aws.String(msg.Key),
				Detail:     aws.String(string(msg.Body)),
				Source:     aws.String(pub.source),
				Time:       aws.Time(time.Now()),
			},
		},
	})
	for _, entry := range resp.Entries {
		if entry.EventId == nil {
			log.Println("injection failed: " + entry.String())
		}
	}
	if err != nil {
		return err
	}
	return nil
}
