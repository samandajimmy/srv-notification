package pubsub

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
)

type HandlerFn = func(ctx context.Context, payload message.Payload) (ack bool, err error)

type Worker struct {
	Topic   string
	Context context.Context
	// Private members
	messages  <-chan *message.Message
	handlerFn HandlerFn
}

func NewWorker(sub message.Subscriber, topic string, args ...interface{}) *Worker {
	// Evaluate arguments
	var ctx context.Context
	if len(args) > 0 {
		arg, ok := args[0].(context.Context)
		if ok {
			ctx = arg
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// Subscribe
	messages, err := sub.Subscribe(ctx, topic)
	if err != nil {
		logger.Errorf("failed to Subscribe. Topic = %s", topic)
	}

	return &Worker{
		Topic:     topic,
		Context:   ctx,
		messages:  messages,
		handlerFn: nil,
	}
}

func (s *Worker) Register(fn HandlerFn) {
	logger.Debugf("registering pubsub function for topic = %s", s.Topic)
	s.handlerFn = fn
}

func (s *Worker) Listen() {
	if s.handlerFn == nil {
		panic(fmt.Errorf("HandlerFn is not initiated for subscriber. Topic = %s", s.Topic))
	}

	for msg := range s.messages {
		logger.Debug("received message. Topic: %s, MessageId: %s", s.Topic, msg.UUID)
		// Call handler
		ack, err := s.handlerFn(msg.Context(), msg.Payload)
		if err != nil {
			logger.Errorf("an error occurred while listening to topic %s. Error = %s", s.Topic, err)
		}

		// If not ack, then retry
		if !ack {
			msg.Nack()
			continue
		}

		// Send done
		msg.Ack()
	}
}
