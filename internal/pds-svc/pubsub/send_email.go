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

type SendEmailHandler struct {
	*Worker
	EmailService contract.EmailService
}

func NewSendEmailHandler(sub message.Subscriber, emailService contract.EmailService) *SendEmailHandler {
	// Init Send Email Handler
	h := SendEmailHandler{
		Worker:       NewWorker(sub, constant.SendEmailTopic),
		EmailService: emailService,
	}

	// Register handler function
	h.Register(h.sendEmail)

	return &h
}

func (h *SendEmailHandler) sendEmail(_ context.Context, payload message.Payload) (ack bool, err error) {
	// Parse payload
	var p dto.SendEmail
	err = json.Unmarshal(payload, &p)
	if err != nil {
		logger.Errorf("failed to parse payload. Topic = %s", h.Topic)
		return false, err
	}

	// Send email
	err = h.EmailService.SendEmail(p)
	if err != nil {
		logger.Errorf("Error when sending email in service %v", err)
		return false, ncore.TraceError(err)
	}

	return true, nil
}
