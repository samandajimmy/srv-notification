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
	var p dto.SendNotificationOptionsRequest
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Errorf("failed to parse payload. Topic = %s", logger.Format(h.Topic), logger.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdKey, p.RequestId)

	// Get service context
	svc := h.svc.WithContext(ctx)

	// Set application
	application := p.Auth

	// Prepare payload push notification
	payloadSendPushNotification := dto.SendPushNotification{
		ApplicationId: application.ID,
		Title:         p.Options.FCM.Title,
		Body:          p.Options.FCM.Body,
		ImageURL:      p.Options.FCM.ImageUrl,
		Token:         p.Options.FCM.Token,
	}
	if p.Options.FCM.Data != nil {
		payloadSendPushNotification.Data = p.Options.FCM.Data
	}

	optionsWebhook := dto.WebhookOptions{
		WebhookURL:       p.Auth.WebhookURL,
		NotificationType: constant.NotificationFCM,
		Notification:     p.Notification,
		Payload:          payloadSendPushNotification,
	}

	log.Debugf("Send to Webhook url: '%s' . Topic :%s", optionsWebhook.WebhookURL, constant.SendFcmTopic)

	// Send Push Notification
	err = svc.SendPushNotificationByTarget(payloadSendPushNotification)
	if err != nil {
		log.Error("Error when sending email in service %v", logger.Error(err), logger.Context(ctx))
		optionsWebhook.NotificationStatus = constant.NotificationStatusFailed
		if optionsWebhook.WebhookURL != "" {
			SendWebhook(optionsWebhook)
		}
		return true, ncore.TraceError(err)
	} else {
		optionsWebhook.NotificationStatus = constant.NotificationStatusSuccess
		if optionsWebhook.WebhookURL != "" {
			SendWebhook(optionsWebhook)
		}
	}

	return true, nil
}
