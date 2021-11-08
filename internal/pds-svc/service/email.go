package service

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nlogger"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nmail"
	"fmt"
	"os"
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
