package handler

import (
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
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
	// Get Payload
	var payload dto.ClientConfig
	err := rx.ParseJSONBody(&payload)
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

	// Call service
	svc := h.Service.WithContext(rx.Context())

	resp, err := svc.CreateClientConfig(payload)
	if err != nil {
		log.Error("failed to create client config.", nlogger.Error(err))
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
	return nil, nil // TODO Implement
}

func (h *ClientConfig) UpdateClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	return nil, nil // TODO Implement
}

func (h *ClientConfig) DeleteClientConfig(rx *nhttp.Request) (*nhttp.Response, error) {
	return nil, nil // TODO Implement
}
