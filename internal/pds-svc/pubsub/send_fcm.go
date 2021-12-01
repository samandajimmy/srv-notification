package pubsub

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type SendFcmPushHandler struct {
	*SubscriberHandler
	NotificationService contract.NotificationService
}

func NewSendFcmPushHandler(sub message.Subscriber, notificationSvc contract.NotificationService) *SendFcmPushHandler {
	// Init Send Email Handler
	h := SendFcmPushHandler{
		SubscriberHandler:   NewSubscriberHandler(sub, constant.SendFcmTopic),
		NotificationService: notificationSvc,
	}

	// Register handler function
	h.Register(h.sendFcm)

	return &h
}

func (h *SendFcmPushHandler) sendFcm(_ context.Context, payload message.Payload) (ack bool, err error) {
	// Parse payload
	var p dto.NotificationCreate
	err = json.Unmarshal(payload, &p)
	if err != nil {
		logger.Errorf("failed to parse payload. Topic = %s", h.Topic)
		return false, err
	}

	// Send email
	err = h.NotificationService.SendNotificationByToken(p)
	if err != nil {
		logger.Errorf("Error when sending email in service %v", err)
		return false, ncore.TraceError(err)
	}

	return true, nil
}
