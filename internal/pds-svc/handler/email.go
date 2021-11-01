package handler

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nhttp"
)

func NewEmail(eventLogService contract.EmailService) *Email {
	return &Email{eventLogService}
}

type Email struct {
	emailService contract.EmailService
}

func (h *Email) PostEmail(rx *nhttp.Request) (*nhttp.Response, error) {

	// Get Payload
	var payload dto.SendEmail
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	err = h.emailService.SendEmail(payload)
	if err != nil {
		log.Errorf("Error when sending email in service %v", err)
		return nil, ncore.TraceError(err)
	}

	return nhttp.OK(), nil
}
