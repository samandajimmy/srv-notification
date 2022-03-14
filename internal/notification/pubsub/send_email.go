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

type SendEmailHandler struct {
	*Worker
	Service *contract.Service
}

func NewSendEmailHandler(sub message.Subscriber, svc *contract.Service) *SendEmailHandler {
	// Init Send Email Handler
	h := SendEmailHandler{
		Worker:  NewWorker(sub, constant.SendEmailTopic),
		Service: svc,
	}

	// Register handler function
	h.Register(h.sendEmail)

	return &h
}

func (h *SendEmailHandler) sendEmail(ctx context.Context, payload message.Payload) (ack bool, err error) {
	// Parse payload
	var p dto.SendNotificationOptionsRequest
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Error("failed to parse payload. Topic = %s", logOption.Format(h.Topic), logOption.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdContextKey, p.RequestId)

	// Get service context
	svc := h.Service.WithContext(ctx)
	defer svc.Close()

	// Set application
	application := p.Auth

	// Set payload send email
	payloadSendEmail := dto.SendEmail{
		ApplicationId: application.ID,
		Subject:       p.Options.SMTP.Subject,
		From: dto.FromFormat{
			Name:  p.Options.SMTP.From.Name,
			Email: p.Options.SMTP.From.Email,
		},
		To:         p.Options.SMTP.To,
		Message:    p.Options.SMTP.Message,
		Attachment: p.Options.SMTP.Attachment,
		MimeType:   p.Options.SMTP.MimeType,
	}

	optionsWebhook := dto.WebhookEmailOptions{
		WebhookURL:       p.Auth.WebhookURL,
		NotificationType: constant.NotificationEmail,
		ApplicationID:    p.Auth.ID,
		Payload:          payloadSendEmail,
	}

	log.Debugf("Send to Webhook url: '%s' . Topic :%s", optionsWebhook.WebhookURL, constant.SendEmailTopic)

	// Send email
	err = svc.SendEmail(&payloadSendEmail)
	if err != nil {
		log.Error("Error when sending email in service", logOption.Error(err))
		optionsWebhook.NotificationStatus = constant.NotificationStatusFailed
		if optionsWebhook.WebhookURL != "" {
			SendWebhookEmail(optionsWebhook)
		}

		return true, ncore.TraceError(err)
	}

	// Send Webhook if email sent
	optionsWebhook.NotificationStatus = constant.NotificationStatusSuccess
	if optionsWebhook.WebhookURL != "" {
		SendWebhookEmail(optionsWebhook)
	}

	return true, nil
}

func SendWebhookEmail(options dto.WebhookEmailOptions) {
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
		"notificationTime":   time.Now().Unix(),
		"applicationId":      options.ApplicationID,
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
