package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func NewEmail(emailService contract.EmailService) *Email {
	return &Email{emailService}
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

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error appear when validate payload %v", err)
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
