package eventbus

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/google/uuid"
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
	if msg == nil {
		return ErrNilMessage
	}
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	d, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := pub.clt.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(pub.eventBusName),
				DetailType:   aws.String(msg.Key),
				Detail:       aws.String(string(d)),
				Source:       aws.String(pub.source),
				Time:         aws.Time(time.Now()),
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
