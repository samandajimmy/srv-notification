package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorilla/mux"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
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
	// TODO: Refactor to Auth Middleware
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

	// TODO: Refactor to Auth Middleware
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
	// Get context
	ctx := rx.Context()

	// Get app
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, ncore.TraceError(err)
	}

	// Get Payload
	var payload dto.SendNotificationOptionsRequest
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Error("Error when parse json body from request", logger.Error(err), logger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("Error when validate payload send notification", logger.Error(err), logger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Normalize request value
	// -- Set user id on fcm options
	if o := payload.Options.FCM; o != nil {
		o.UserId = payload.UserId
	}

	if o := payload.Options.SMTP; o != nil {
		o.UserId = payload.UserId
	}

	// -- Set app to auth
	payload.Auth = app

	// -- Set request id
	payload.RequestId = GetRequestId(rx)

	// Create notification
	svc := h.Service.WithContext(rx.Context())
	data, err := svc.CreateNotification(payload)
	if err != nil {
		log.Error("Error when create notification", logger.Error(err), logger.Context(rx.Context()))
		return nil, err
	}

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

	return nhttp.Success().SetData(data), nil
}

func (h *Notification) GetDetailNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get id
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.GetNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id

	err := payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.GetDetailNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) GetDetailNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get id
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.GetNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id

	err := payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.GetDetailNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) DeleteNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get ID
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.GetNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id

	err := payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	err = svc.DeleteNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}
