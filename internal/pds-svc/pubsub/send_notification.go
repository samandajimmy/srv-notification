package pubsub

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/net/context"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type SendNotificationHandler struct {
	*Worker
	svc *contract.Service
}

func NewSendNotificationHandler(sub message.Subscriber, notificationSvc *contract.Service) *SendNotificationHandler {
	// Init Send Notification Handler
	h := SendNotificationHandler{
		Worker: NewWorker(sub, constant.SendNotificationTopic),
		svc:    notificationSvc,
	}

	// Register handler function
	h.Register(h.sendNotification)

	return &h
}

func (h *SendNotificationHandler) sendNotification(ctx context.Context, payload message.Payload) (ack bool, err error) {
	// Parse payload
	var p dto.SendNotificationOptionsRequest
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Error("failed to parse payload. Topic = %s", logger.Format(h.Topic), logger.Error(err))
		return false, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdKey, p.RequestId)

	// Get service context
	svc := h.svc.WithContext(ctx)

	// persist to DB
	err = svc.CreateNotification(p)
	if err != nil {
		return false, ncore.TraceError(err)
	}

	// TODO Send Email

	// TODO Send Push Notification

	return true, nil
}
