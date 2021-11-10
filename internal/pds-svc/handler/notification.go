package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
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
