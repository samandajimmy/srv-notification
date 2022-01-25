package pubsub

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/net/context"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
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
		return true, err
	}

	// Get service context
	svc := h.svc.WithContext(ctx)

	application := p.Auth

	// decode html
	// Send email
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
	// Send to Email
	err = svc.SendEmail(payloadSendEmail)
	if err != nil {
		log.Error("failed while sending email. Topic = %s. Err %v.", logger.Format(h.Topic), logger.Error(err), err)
		return true, err
	}

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
	// send to fcm
	err = svc.SendPushNotificationByTarget(payloadSendPushNotification)
	if err != nil {
		log.Error("Error when sending email in service %v", logger.Error(err), logger.Context(ctx))
		return true, err
	}

	//svc.CreateNotification()

	return true, nil
}
