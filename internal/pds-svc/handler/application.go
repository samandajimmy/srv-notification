package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type Application struct {
	Service *contract.Service
}

func NewApplication(svc *contract.Service) *Application {
	return &Application{
		Service: svc,
	}
}

func (h *Application) PostApplication(rx *nhttp.Request) (*nhttp.Response, error) {
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
