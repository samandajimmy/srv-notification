package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gorilla/mux"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/logger"
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
	payload.Notification = data
	pubsubPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: unable to marshal payload")
	}
	// Init massage payload for pubsub
	msg := message.NewMessage(watermill.NewUUID(), pubsubPayload)

	// Publish to Email
	err = h.publisher.Publish(constant.SendEmailTopic, msg)
	if err != nil {
		log.Errorf("failed to publish message to topic = %s", constant.SendEmailTopic)
		return nil, err
	}

	// Publish to Fcm
	err = h.publisher.Publish(constant.SendFcmTopic, msg)
	if err != nil {
		log.Errorf("failed to publish message to topic = %s", constant.SendFcmTopic)
		return nil, err
	}

	return nhttp.Success().SetData(data), nil
}

func (h *Notification) GetDetailNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logger.Format(err), logger.Error(err), logger.Context(ctx))
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
		log.Error("error when validate payload", logger.Error(err), logger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.GetDetailNotification(payload)
	if err != nil {
		log.Error("error when call service", logger.Error(err), logger.Context(ctx))
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
		log.Errorf("Error when get application: %s", logger.Format(err), logger.Error(err), logger.Context(ctx))
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
	err = svc.DeleteNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}

func (h *Notification) GetCountNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logger.Format(err), logger.Error(err), logger.Context(ctx))
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
		log.Error("request validate error", logger.Error(err), logger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	resp, err := svc.GetCountNotification(payload)
	if err != nil {
		log.Error("error when call service err: %v", logger.Error(err), logger.Context(ctx))
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) GetListNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get context
	ctx := rx.Context()

	// Get application
	app, err := getApplication(rx)
	if err != nil {
		log.Errorf("Error when get application: %s", logger.Format(err), logger.Error(err), logger.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get list parameters
	q := rx.URL.Query()
	//Get parameter
	listParam := dto.NotificationFindOptions{
		FindOptions: dto.FindOptions{
			Limit:         nval.ParseIntFallback(q.Get("limit"), 10),
			Skip:          nval.ParseIntFallback(q.Get("skip"), 0),
			SortBy:        nval.ParseStringFallback(q.Get("sortBy"), "createdAt"),
			SortDirection: nval.ParseStringFallback(q.Get("sortDirection"), "desc"),
			Filters:       map[string]interface{}{},
		},
	}

	//Call service
	svc := h.Service.WithContext(rx.Context())
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	if app != nil {
		listParam.Filters["applicationId"] = app.ID
	}

	if v := nval.ParseInt64Fallback(q.Get("filters[userId]"), 0); v > 0 {
		listParam.Filters["userId"] = v
	} else {
		log.Error("userId is required to get list data", nlogger.Error(err))
		return nil, nhttp.BadRequestError
	}

	resp, err := svc.ListNotification(listParam)
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
		log.Errorf("Error when get application: %s", logger.Format(err), logger.Error(err), logger.Context(ctx))
		return nil, ncore.TraceError(err)
	}

	// Get user id
	id := mux.Vars(rx.Request)["id"]

	// Set payload
	var payload dto.UpdateIsReadNotification
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
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err), logger.Context(ctx))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	resp, err := svc.UpdateIsRead(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
