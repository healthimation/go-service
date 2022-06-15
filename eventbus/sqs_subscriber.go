package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
	"time"
)

type sqsSubscriber struct {
	clt       *sqs.SQS
	timeout   time.Duration
	maxWorker int
	maxMsg    int
	queueUrl  string

	mu             sync.Mutex
	handlers       map[string]MessageHandler
	keys           []string
	defaultHandler MessageHandler
	exit           chan bool
}

type SqsMessage struct {
	Version    string    `json:"version"`
	ID         string    `json:"id"`
	DetailType string    `json:"detail-type"`
	Source     string    `json:"source"`
	Account    string    `json:"account"`
	Time       time.Time `json:"time"`
	Region     string    `json:"region"`
	Detail     *Message  `json:"detail"`
}

func defaultHandler(message *Message) {
	log.Println("default handler", message)
}

func NewSQSSubscriber(sess *session.Session, cfg *SubscriberConfig) (Subscriber, error) {
	svc := sqs.New(sess)
	if cfg.DefaultHandler == nil {
		cfg.DefaultHandler = defaultHandler
	}
	return &sqsSubscriber{
		clt:            svc,
		queueUrl:       cfg.QueueUrl,
		maxMsg:         cfg.MaxMsg,
		maxWorker:      cfg.MaxWorker,
		handlers:       map[string]MessageHandler{},
		defaultHandler: cfg.DefaultHandler,
		exit:           make(chan bool),
		timeout:        cfg.Timeout,
	}, nil
}

func (c *sqsSubscriber) Subscribe(key string, fn MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[key] = fn
}

func (c *sqsSubscriber) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(c.maxWorker)

	for i := 1; i <= c.maxWorker; i++ {
		go c.worker(ctx, wg, i)
	}

	wg.Wait()

	<-c.exit
	return nil
}

func (c *sqsSubscriber) Stop() {
	c.exit <- true
}

func (c *sqsSubscriber) getQueueUrl(key string) string {
	return key
}

func (c *sqsSubscriber) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	log.Printf("worker %d: started\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d: stopped\n", id)
			return
		default:
		}

		log.Println("Waiting for messages")
		msgs, err := c.receive(ctx, c.queueUrl, int64(c.maxMsg))
		if err != nil {
			// Critical error!
			log.Printf("worker %d: receive error: %s\n", id, err.Error())
			continue
		}

		if len(msgs) == 0 {
			continue
		}

		log.Printf("Received %d messages\n", len(msgs))
		msgWg := &sync.WaitGroup{}
		msgWg.Add(len(msgs))

		for _, m1 := range msgs {
			go func(msg *sqs.Message, w *sync.WaitGroup) {
				defer w.Done()
				defer func() {
					err := c.delete(ctx, c.queueUrl, *msg.ReceiptHandle)
					if err != nil {
						log.Printf("worker %d: delete error: %s\n", id, err.Error())
					}
				}()
				sqsMsg := &SqsMessage{}
				err = json.Unmarshal([]byte(*msg.Body), sqsMsg)
				if err != nil {
					log.Printf("worker %d: json unmarshal sqsMessage error: %s\n", id, err.Error())
					return
				}
				fn, ok := c.handlers[sqsMsg.DetailType]
				if !ok {
					c.defaultHandler(sqsMsg.Detail)
					return
				}
				fn(sqsMsg.Detail)
			}(m1, msgWg)
		}
		log.Println("Waiting for messages to be processed")
		msgWg.Wait()
		log.Println("Messages are processed")

		//if c.config.Type == SyncConsumer {
		//	c.sync(ctx, msgs)
		//} else {
		//	c.async(ctx, msgs)
		//}
	}
}

//func (c *sqsSubscriber) sync(ctx context.Context, msgs []*sqs.Message) {
//	for _, msg := range msgs {
//		c.consume(ctx, msg)
//	}
//}
//
//func (c *sqsSubscriber) async(ctx context.Context, msgs []*sqs.Message) {
//	wg := &sync.WaitGroup{}
//	wg.Add(len(msgs))
//
//	for _, msg := range msgs {
//		go func(msg *sqs.Message) {
//			defer wg.Done()
//
//			c.consume(ctx, msg)
//		}(msg)
//	}
//
//	wg.Wait()
//}

func (c *sqsSubscriber) receive(ctx context.Context, queueURL string, maxMsg int64) ([]*sqs.Message, error) {
	if maxMsg < 1 || maxMsg > 10 {
		return nil, fmt.Errorf("receive argument: msgMax valid values: 1 to 10: given %d", maxMsg)
	}

	var waitTimeSeconds int64 = 10

	// Must always be above `WaitTimeSeconds` otherwise `ReceiveMessageWithContext`
	// trigger context timeout error.
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(waitTimeSeconds+5))
	defer cancel()

	res, err := c.clt.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   aws.Int64(maxMsg),
		WaitTimeSeconds:       aws.Int64(waitTimeSeconds),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	return res.Messages, nil
}

func (c *sqsSubscriber) delete(ctx context.Context, queueURL, rcvHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if _, err := c.clt.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
