package events

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Queue interface {
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Publish(topic string, messages ...*message.Message) error
}
