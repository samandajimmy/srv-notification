package handler

import (
	"encoding/json"
	"errors"
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
	payload.Notification = data
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

func (h *Notification) GetCountNotification(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get user id
	userId := rx.URL.Query().Get("userId")
	if nval.ParseInt64Fallback(userId, 0) <= 0 {
		return nil, nhttp.BadRequestError.Wrap(errors.New("User id is required."))
	}

	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.BadRequestError
	}

	// Set payload
	var payload dto.GetCountNotification
	payload.RequestId = GetRequestId(rx)
	payload.UserRefId = nval.ParseInt64Fallback(userId, 0)

	// Call service
	svc := h.Service.WithContext(rx.Context())
	application, err := svc.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}
	payload.Application = application

	resp, err := svc.GetCountNotification(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Notification) GetListNotification(rx *nhttp.Request) (*nhttp.Response, error) {
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

	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.BadRequestError
	}

	//Call service
	svc := h.Service.WithContext(rx.Context())
	application, err := svc.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	if application != nil {
		listParam.Filters["applicationId"] = application.ID
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
	// Get user id
	id := mux.Vars(rx.Request)["id"]

	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.BadRequestError
	}

	// Set payload
	var payload dto.GetCountNotification
	payload.RequestId = GetRequestId(rx)
	payload.ID = id

	err := payload.Validate()
	if err != nil {
		log.Errorf("id is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	application, err := svc.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to auth application", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}
	payload.Application = application

	resp, err := svc.UpdateIsRead(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}
