package pubsub

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	logOption "github.com/nbs-go/nlogger/v2/option"
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
		log.Error("failed to Subscribe. Topic = %s", logOption.Format(topic), logOption.Error(err))
	}

	return &Worker{
		Topic:     topic,
		Context:   ctx,
		messages:  messages,
		handlerFn: nil,
	}
}

func (s *Worker) Register(fn HandlerFn) {
	log.Debugf("registering pubsub function for topic = %s", s.Topic)
	s.handlerFn = fn
}

func (s *Worker) Listen() {
	if s.handlerFn == nil {
		panic(fmt.Errorf("HandlerFn is not initiated for subscriber. Topic = %s", s.Topic))
	}

	for msg := range s.messages {
		log.Debugf("received message. Topic: %s, MessageId: %s", s.Topic, msg.UUID)
		// Call handler
		ack, err := s.handlerFn(msg.Context(), msg.Payload)
		if err != nil {
			log.Error("an error occurred while listening to topic %s", logOption.Format(s.Topic), logOption.Error(err))
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
