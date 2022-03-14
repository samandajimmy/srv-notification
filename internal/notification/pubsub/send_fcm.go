package pubsub

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"

	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"time"
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
		log.Errorf("failed to parse payload. Topic = %s", logOption.Format(h.Topic), logOption.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdContextKey, p.RequestId)

	// Get service context
	svc := h.svc.WithContext(ctx)
	defer svc.Close()

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
	err = svc.SendPushNotificationByTarget(&payloadSendPushNotification)
	if err != nil {
		log.Error("Error when sending email in service %v", logOption.Error(err), logOption.Context(ctx))
		optionsWebhook.NotificationStatus = constant.NotificationStatusFailed
		if optionsWebhook.WebhookURL != "" {
			SendWebhookFcm(optionsWebhook)
		}
		return true, ncore.TraceError(err)
	}

	// Send webhook when fcm sent
	optionsWebhook.NotificationStatus = constant.NotificationStatusSuccess
	if optionsWebhook.WebhookURL != "" {
		SendWebhookFcm(optionsWebhook)
	}

	return true, nil
}

func SendWebhookFcm(options dto.WebhookOptions) {
	// Set header
	reqHeader := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Set payload
	reqBody := map[string]interface{}{
		"notificationId":     "N/A", // N/A
		"userId":             "N/A", // N/A
		"notificationStatus": options.NotificationStatus,
		"notificationTime":   time.Now(),
		"applicationId":      options.Notification.ApplicationId,
		"payload":            options.Payload,
	}

	// Initiate client
	c := nclient.NewNucleoClient(options.WebhookURL)

	// Send webhook to client
	_, err := c.PostData(reqHeader, reqBody)
	if err != nil {
		log.Error("error when send webhook to client", logOption.Error(err))
	}
}
