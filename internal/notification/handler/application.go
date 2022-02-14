package handler

import (
	"errors"
	"github.com/gorilla/mux"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

type Application struct {
	Service *contract.Service
}

func NewApplication(svc *contract.Service) *Application {
	return &Application{
		Service: svc,
	}
}

func (h *Application) PostCreateApplication(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get authenticated entity
	subject, err := GetSubject(rx)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Get Payload
	var payload dto.Application
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error appear when validate payload: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set request id
	payload.RequestId = GetRequestId(rx)
	payload.Subject = subject

	// Call service
	svc := h.Service.WithContext(rx.Context())

	resp, err := svc.CreateApplication(payload)
	if err != nil {
		return nil, err
	}

	// Set payload
	return nhttp.Success().SetData(resp), nil
}

func (h *Application) GetFindApplication(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get list parameters
	q := rx.URL.Query()
	//Get parameter
	listParam := dto.ApplicationFindOptions{
		FindOptions: dto.FindOptions{
			Limit:         nval.ParseIntFallback(q.Get("limit"), 10),
			Skip:          nval.ParseIntFallback(q.Get("skip"), 0),
			SortBy:        nval.ParseStringFallback(q.Get("sortBy"), "createdAt"),
			SortDirection: nval.ParseStringFallback(q.Get("sortDirection"), "desc"),
			Filters:       map[string]interface{}{},
		},
	}

	if v := q.Get("xid"); v != "" {
		listParam.Filters["xid"] = v
	}

	if v := q.Get("name"); v != "" {
		listParam.Filters["name"] = v
	}

	//Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.ListApplication(&listParam)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK().SetData(resp), nil
}

func (h *Application) GetDetailApplication(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.GetApplication
	payload.RequestId = GetRequestId(rx)
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.GetApplication(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Application) PutUpdateApplication(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Auth Subject
	actor, err := GetSubject(rx)
	if err != nil {
		return nil, err
	}

	// Get Payload
	var payload dto.ApplicationUpdateOptions
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set modifier and id
	payload.RequestId = GetRequestId(rx)
	payload.XID = mux.Vars(rx.Request)["xid"]
	payload.Subject = actor

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())

	resp, err := svc.UpdateApplication(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *Application) DeleteApplication(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.GetApplication
	payload.RequestId = GetRequestId(rx)
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	err := svc.DeleteApplication(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}
