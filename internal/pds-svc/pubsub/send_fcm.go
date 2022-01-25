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
	var p dto.SendPushNotificationRequest
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Error("failed to parse payload. Topic = %s", logger.Format(h.Topic), logger.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdKey, p.RequestId)

	// Get service context
	svc := h.svc.WithContext(ctx)

	pushNotificationPayload := dto.SendPushNotification{
		RequestId:     p.RequestId,
		Title:         p.Title,
		Body:          p.Body,
		ImageURL:      p.ImageUrl,
		Token:         p.Token,
		ApplicationId: p.Auth.ID,
		Data:          p.Data,
	}

	// Send email
	err = svc.SendPushNotificationByTarget(pushNotificationPayload)
	if err != nil {
		log.Error("Error when sending email in service %v", logger.Error(err), logger.Context(ctx))
		return true, ncore.TraceError(err)
	}

	// prepare FCM Options
	fcmOption := &dto.FCMOption{
		UserId:   p.UserId,
		Title:    p.Title,
		Body:     p.Body,
		ImageUrl: p.ImageUrl,
		Token:    p.Token,
		Metadata: p.Metadata,
		Data:     p.Data,
	}
	// prepare create notification
	notification := dto.SendNotificationOptionsRequest{
		UserId:    p.UserId,
		RequestId: p.RequestId,
		Auth:      p.Auth,
		Options: dto.NotificationOptionVO{
			FCM: fcmOption,
		},
	}
	// execute Create Notification service
	err = svc.CreateNotification(notification)
	if err != nil {
		log.Error("Error when create notification %v", logger.Error(err), logger.Context(ctx))
		return true, err
	}

	return true, nil
}
