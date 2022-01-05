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
	var p dto.SendEmail
	err = json.Unmarshal(payload, &p)
	if err != nil {
		log.Error("failed to parse payload. Topic = %s", logger.Format(h.Topic), logger.Error(err))
		return true, err
	}

	// Set request id to context
	ctx = context.WithValue(ctx, nhttp.RequestIdKey, p.RequestId)

	// Get service context
	svc := h.Service.WithContext(ctx)

	// Send email
	err = svc.SendEmail(p)
	if err != nil {
		log.Error("Error when sending email in service", logger.Error(err))
		return true, ncore.TraceError(err)
	}

	return true, nil
}
