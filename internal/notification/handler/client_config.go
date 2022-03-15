package handler

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
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
	subject := GetSubject(rx)

	// Get Payload
	var payload dto.ClientConfigRequest
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request.", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	log.Debugf("Received Create ClientConfig request.")
	err = payload.Validate()
	if err != nil {
		log.Error("Error appear when validate payload.", logOption.Error(err))
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set request id
	payload.RequestId = rx.GetRequestId()
	payload.Subject = subject

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.CreateClientConfig(&payload)
	if err != nil {
		log.Error("failed to create client.", logOption.Error(err))
		return nil, err
	}

	// Set payload
	return nhttp.Success().SetData(resp), nil
}

func (h *ClientConfig) ListClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	payload, err := getListPayload(rx)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Call service
	srv := h.Service.WithContext(rx.Context())
	respData, err := srv.ListClientConfig(payload)
	if err != nil {
		log.Error("failed to call service.", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set response
	resp := nhttp.OK().SetData(respData)

	return resp, nil
}

func (h *ClientConfig) GetDetailClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get xid
	xid := mux.Vars(rx.Request)["xid"]
	if xid == "" {
		err := errors.New("xid is not found on params")
		log.Errorf("xid is not found on params. err: %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	var payload dto.ClientConfigRequest
	payload.RequestId = rx.GetRequestId()
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.GetDetailClientConfig(&payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.Success().SetData(resp), nil
}

func (h *ClientConfig) UpdateClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	// Get Auth Subject
	subject := GetSubject(rx)

	// Get Payload
	var payload dto.ClientConfigUpdateOptions
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set modifier and id
	payload.RequestId = rx.GetRequestId()
	payload.XID = mux.Vars(rx.Request)["xid"]
	payload.Subject = subject

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	resp, err := svc.UpdateClientConfig(&payload)
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
	payload.RequestId = rx.GetRequestId()
	payload.XID = xid

	// Call service
	svc := h.Service.WithContext(rx.Context())
	defer svc.Close()

	err := svc.DeleteClientConfig(&payload)
	if err != nil {
		log.Errorf("error when call service err: %v", err)
		return nil, err
	}

	return nhttp.OK(), nil
}
