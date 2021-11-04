package handler

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nhttp"
)

func NewNotification(notificationService contract.NotificationService) *Notification {
	return &Notification{notificationService}
}

type Notification struct {
	notificationService contract.NotificationService
}

func (h *Notification) PostNotification(rx *nhttp.Request) (*nhttp.Response, error) {

	// Get Payload
	var payload dto.NotificationCreate
	err := rx.ParseJSONBody(&payload)
	if err != nil {
		log.Errorf("Error when parse json body from request %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Validate payload
	err = payload.Validate()
	if err != nil {
		log.Errorf("Error when validate payload notification %v", err)
		return nil, nhttp.BadRequestError.Wrap(err)
	}

	// Set payload
	err = h.notificationService.SendNotificationByToken(payload)
	if err != nil {
		log.Errorf("Error when send notification by token %v", err)
		return nil, ncore.TraceError(err)
	}

	return nhttp.OK(), nil
}
