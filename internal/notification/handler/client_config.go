package handler

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

type ClientConfig struct {
	Service *contract.Service
}

func NewClientConfig(svc *contract.Service) *ClientConfig {
	return &ClientConfig{
		Service: svc,
	}
}

func (h *ClientConfig) CreateClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get authenticated entity
	subject, err := GetSubject(rx)
	if err != nil {
		log.Errorf("Error when get subject authenticated entity.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}
	// Get Payload
	var payload dto.ClientConfigRequest
	err = rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request.", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received Create ClientConfig request.")
	err = payload.Validate()
	if err != nil {
		log.Error("Error appear when validate payload.", nlogger.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set request id
	payload.RequestId = GetRequestId(rx)
	payload.Subject = subject

	// Call service
	svc := h.Service.WithContext(rx.Context())

	resp, err := svc.CreateClientConfig(payload)
	if err != nil {
		log.Error("failed to create client .", nlogger.Error(err))
		return nil, err
	}

	// Set payload
	return nhttp.Success().SetData(resp), nil
}

func (h *ClientConfig) SearchClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get authenticated entity
	subject, err := GetSubject(rx)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Get list parameters
	q := rx.URL.Query()
	// Get parameter
	listParam := dto.ClientConfigFindOptions{
		FindOptions: dto.FindOptions{
			Limit:         nval.ParseIntFallback(q.Get("limit"), 10),
			Skip:          nval.ParseIntFallback(q.Get("skip"), 0),
			SortBy:        nval.ParseStringFallback(q.Get("sortBy"), "createdAt"),
			SortDirection: nval.ParseStringFallback(q.Get("sortDirection"), "desc"),
			Filters:       map[string]interface{}{},
		},
		Subject: subject,
	}

	if v := q.Get("applicationXid"); v != "" {
		listParam.Filters["applicationXid"] = v
	}

	// Call service
	srv := h.Service.WithContext(rx.Context())
	respData, err := srv.ListClientConfig(listParam)
	if err != nil {
		log.Error("failed to call service.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	// Set response
	resp := nhttp.OK().SetData(respData)

	return resp, nil
}

func (h *ClientConfig) DetailClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.ClientConfigRequest
	payload.RequestId = GetRequestId(rx)
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	resp, err := svc.GetClientConfig(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *ClientConfig) UpdateClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Auth Subject
	actor, err := GetSubject(rx)
	if err != nil {
		return nil, err
	}

	// Get Payload
	var payload dto.ClientConfigUpdateOptions
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

	resp, err := svc.UpdateClientConfig(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *ClientConfig) DeleteClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.GetClientConfig
	payload.RequestId = GetRequestId(rx)
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	err := svc.DeleteClientConfig(payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}
