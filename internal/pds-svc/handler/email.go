package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func NewEmail(emailService contract.EmailService, publisher message.Publisher) *Email {
	return &Email{emailService, publisher}
}

type Email struct {
	emailService contract.EmailService
	publisher    message.Publisher
}

func (h *Email) PostEmail(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Payload
	var payload dto.SendEmail
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received PostEmail request. Payload: %+v", payload)
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error appear when validate payload: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Publish to pubsub
	pubsubPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: unable to marshal payload")
	}

	msg := message.NewMessage(watermill.NewUUID(), pubsubPayload)
	err = h.publisher.Publish(constant.SendEmailTopic, msg)
	if err != nil {
		log.Errorf("failed to publish message to topic = %s", constant.SendEmailTopic)
		return nil, err
	}

	// Set payload
	return nhttp.OK(), nil
}
