package eventbus

import (
	"context"
	"errors"
)

var (
	ErrNilMessage = errors.New("nil message")
)

type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
}
