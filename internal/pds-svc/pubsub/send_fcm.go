package pubsub

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type SendFcmPushHandler struct {
	*Worker
	svc *contract.Service
}

func NewSendFcmPushHandler(sub message.Subscriber, notificationSvc *contract.Service) *SendFcmPushHandler {
	// Init Send Email Handler
	h := SendFcmPushHandler{
		Worker: NewWorker(sub, constant.SendFcmTopic),
		svc:    notificationSvc,
	}

	// Register handler function
	h.Register(h.sendFcm)

	return &h
}

func (h *SendFcmPushHandler) sendFcm(ctx context.Context, payload message.Payload) (ack bool, err error) {
	// Parse payload
	var p dto.SendPushNotification
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Error("failed to parse payload. Topic = %s", logger.Format(h.Topic), logger.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdKey, p.RequestId)

	// Get service context
	svc := h.svc.WithContext(ctx)

	// Send email
	err = svc.SendPushNotificationByTarget(p)
	if err != nil {
		log.Error("Error when sending email in service %v", logger.Error(err), logger.Context(ctx))
		return true, ncore.TraceError(err)
	}

	return true, nil
}
