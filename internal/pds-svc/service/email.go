package service

import (
	"fmt"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nmail"
)

var log = nlogger.Get()

type Email struct {
	mailer    *nmail.NucleoEmail
	responses *ncore.ResponseMap
}

func (s *Email) HasInitialized() bool {
	return true
}

func (s *Email) Init(app *contract.PdsApp) error {

	var err error

	s.responses = app.Responses

	s.mailer, err = nmail.NewNucleoEmail(app.Config.SMTP)
	if err != nil {
		log.Error("Error when initialize nmail! %v", err)
		return ncore.TraceError(err)
	}

	return nil
}

func (s *Email) SendEmail(payload dto.SendEmail) error {

	// Set message
	message, err := s.mailer.ComposeEmail(payload)
	if err != nil {
		log.Error("Error when trying to compose email!", nlogger.Error(err))
		return ncore.TraceError(err)
	}

	err = s.mailer.Client.DialAndSend(message)
	if err != nil {
		log.Error("Failed dial and send email!", nlogger.Error(err))
		return ncore.TraceError(err)
	}

	// Remove file if email has been sent
	if payload.Attachment != "" {
		name := fmt.Sprintf("tmp/tmp.%v", payload.MimeType)
		err = os.Remove(name)
		if err != nil {
			log.Error("Failed remove file after email sent!", nlogger.Error(err))
			return ncore.TraceError(err)
		}
	}

	log.Infof("Email has been sent successfully!.")
	return nil
}
