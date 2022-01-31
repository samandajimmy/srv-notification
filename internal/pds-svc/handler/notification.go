package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func NewNotification(publisher message.Publisher, service *contract.Service) *Notification {
	return &Notification{
		publisher,
		service,
	}
}

type Notification struct {
	publisher message.Publisher
	Service   *contract.Service
}

func (h *Notification) PostNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.BadRequestError
	}

	// Get Payload
	var payload dto.SendPushNotificationRequest
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received PostNotification request.")
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate payload notification %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service Auth
	srv := h.Service.WithContext(rx.Context())
	application, err := srv.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}
	if application != nil {
		payload.Auth = application
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

func (h *Notification) PostCreateNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// TODO: Refactor to Auth Middleware
	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.BadRequestError
	}
	// Get Payload
	var payload dto.SendNotificationOptionsRequest
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received SendNotification request.")
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate payload send notification %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	srv := h.Service.WithContext(rx.Context())
	application, err := srv.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to call service.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}
	if application != nil {
		payload.Auth = application
	}
	// Set request id
	payload.RequestId = GetRequestId(rx)

	// Publish to pubsub
	pubsubPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: unable to marshal payload")
	}

	msg := message.NewMessage(watermill.NewUUID(), pubsubPayload)
	err = h.publisher.Publish(constant.SendNotificationTopic, msg)
	if err != nil {
		log.Errorf("failed to publish message to topic = %s", constant.SendNotificationTopic)
		return nil, err
	}

	return nhttp.OK(), nil
}
