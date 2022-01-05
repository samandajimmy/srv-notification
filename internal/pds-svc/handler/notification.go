package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func NewNotification(publisher message.Publisher) *Notification {
	return &Notification{publisher}
}

type Notification struct {
	publisher message.Publisher
}

func (h *Notification) PostNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.SendPushNotification
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received PostNotification request. Payload: %+v", payload)
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate payload notification %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set request id
	payload.RequestId = GetRequestId(rx)

	// Publish to pubsub
	pubsubPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: unable to marshal payload")
	}

	msg := message.NewMessage(watermill.NewUUID(), pubsubPayload)
	err = h.publisher.Publish(constant.SendFcmTopic, msg)
	if err != nil {
		log.Errorf("failed to publish message to topic = %s", constant.SendFcmTopic)
		return nil, err
	}

	return nhttp.OK(), nil
}
