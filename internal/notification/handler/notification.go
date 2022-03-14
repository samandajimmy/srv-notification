package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorilla/mux"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"

	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
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
		log.Error("Error when parse json body from request", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Error("Error when validate payload send notification", logOption.Error(err), logOption.Context(ctx))
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

	// -- Set subject
	payload.Subject = GetSubject(rx)

	// Create notification
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	switch {
	case payload.Options.SMTP != nil && payload.Options.FCM == nil:
		// SMTP Section
		pubSubPayload, errM := json.Marshal(payload)
		if errM != nil {
			return nil, fmt.Errorf("unexpected error: unable to marshal payload")
		}

		// Set message
		msg := message.NewMessage(watermill.NewUUID(), pubSubPayload)
		err = h.publisher.Publish(constant.SendEmailTopic, msg)
		if err != nil {
			log.Errorf("failed to publish message to topic = %s", constant.SendEmailTopic)
			return nil, err
		}

		return nhttp.OK(), nil
	case payload.Options.FCM != nil && payload.Options.SMTP == nil:
		// FCM Section
		data, errCreate := svc.CreateNotification(&payload)
		if errCreate != nil {
			log.Error("error found when create notification", logOption.Error(errCreate), logOption.Context(ctx))
			return nil, errCreate
		}
		payload.Notification = data
		pubSubPayload, errM := json.Marshal(payload)
		if errM != nil {
			return nil, fmt.Errorf("unexpected error: unable to marshal payload")
		}

		// Set message
		msg := message.NewMessage(watermill.NewUUID(), pubSubPayload)
		// Publish to FCM
		err = h.publisher.Publish(constant.SendFcmTopic, msg)
		if err != nil {
			log.Errorf("failed to publish message to topic = %s", constant.SendFcmTopic)
			return nil, err
		}

		return nhttp.Success().SetData(data), nil
	case payload.Options.FCM != nil && payload.Options.SMTP != nil:
		// FCM & SMTP Section
		data, errCreate := svc.CreateNotification(&payload)
		if errCreate != nil {
			log.Error("Error when create notification", logOption.Error(errCreate), logOption.Context(ctx))
			return nil, errCreate
		}

		payload.Notification = data
		pubSubPayload, errM := json.Marshal(payload)
		if errM != nil {
			return nil, fmt.Errorf("unexpected error: unable to marshal payload")
		}

		// Set message
		msg := message.NewMessage(watermill.NewUUID(), pubSubPayload)

		// Publish to fcm
		err = h.publisher.Publish(constant.SendFcmTopic, msg)
		if err != nil {
			log.Errorf("failed to publish message to topic = %s", constant.SendFcmTopic)
			return nil, err
		}

		// Publish to email
		err = h.publisher.Publish(constant.SendEmailTopic, msg)
		if err != nil {
			log.Errorf("failed to publish message to topic = %s", constant.SendEmailTopic)
			return nil, err
		}

		return nhttp.Success().SetData(data), nil
	}

	return nhttp.BadRequest(nil), nil
}

func (h *Notification) GetDetailNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logOption.Format(err), logOption.Error(err), logOption.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get id
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.GetNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id
	payload.Application = app

	err = payload.Validate()
	if err != nil {
		log.Error("error when validate payload", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.GetDetailNotification(&payload)
	if err != nil {
		log.Error("error when call service", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) DeleteNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logOption.Format(err), logOption.Error(err), logOption.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get ID
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.GetNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id
	payload.Application = app

	err = payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	err = svc.DeleteNotification(&payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *Notification) CountNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logOption.Format(err), logOption.Error(err), logOption.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get user id
	userId := nval.ParseInt64Fallback(rx.URL.Query().Get("userId"), 0)

	// Set payload
	var payload dto.GetCountNotification
	payload.RequestId = GetRequestId(rx)
	payload.UserRefId = userId
	payload.Application = app

	err = payload.Validate()
	if err != nil {
		log.Error("request validate error", logOption.Error(err), logOption.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.CountNotification(&payload)
	if err != nil {
		log.Error("error when call service err: %v", logOption.Error(err), logOption.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) ListNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	payload, err := getListPayload(rx)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Validate payload
	ctx := rx.Context()
	if payload.Filters == nil {
		log.Error("invalid empty filters query in ListNotification", logOption.Context(ctx))
		return nil, nhttp.BadRequestError
	}

	if v, ok := payload.Filters[constant.UserRefIdKey]; !ok || v == "" {
		log.Error("filters[userId] is required", logOption.Context(ctx))
		return nil, nhttp.BadRequestError
	}

	// Set application Id key
	app, err := getApplication(rx)
	if err != nil {
		return nil, ncore.TraceError(err)
	}
	payload.Filters[constant.ApplicationIdKey] = fmt.Sprintf("%d", app.ID)

	svc := h.Service.WithContext(ctx)
	defer svc.Close()

	resp, err := svc.ListNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK().SetData(resp), nil
}

func (h *Notification) UpdateIsReadNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logOption.Format(err), logOption.Error(err), logOption.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get user id
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	payload := dto.UpdateIsReadNotification{
		RequestId:   GetRequestId(rx),
		Application: app,
		ID:          id,
		Subject:     GetSubject(rx),
	}

	err = payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.UpdateIsRead(&payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
