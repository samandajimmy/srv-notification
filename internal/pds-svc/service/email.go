package service

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nlogger"
	"encoding/base64"
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

var log = nlogger.Get()

type Email struct {
	Host      string
	Port      int
	Username  string
	Password  string
	responses *ncore.ResponseMap
}

func (s *Email) HasInitialized() bool {
	return true
}

func (s *Email) Init(app *contract.PdsApp) error {
	s.responses = app.Responses

	smtp := &app.Config.SMTP

	s.Host = smtp.Host
	s.Port = smtp.Port
	s.Username = smtp.Username
	s.Password = smtp.Password

	return nil
}

func (s *Email) SendEmail(payload dto.SendEmail) error {

	mailer, dialer, err := s.ComposeEmail(payload)
	if err != nil {
		log.Error("Error when trying to compose email!", nlogger.Error(err))
		return ncore.TraceError(err)
	}

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Error("Failed dial and send email!", nlogger.Error(err))
		return ncore.TraceError(err)
	}

	// Remove file if email has been sent
	name := fmt.Sprintf("tmp/tmp.%v", payload.MimeType)
	err = os.Remove(name)
	if err != nil {
		log.Error("Failed remove file after email sent!", nlogger.Error(err))
		return ncore.TraceError(err)
	}

	return nil
}

func (s *Email) ComposeEmail(payload dto.SendEmail) (*gomail.Message, *gomail.Dialer, error) {

	// Set dialer
	dialer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// Set message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "no-reply@email.com")
	mailer.SetHeader("To", payload.Subject)
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", payload.Message)

	if payload.Attachment != "" {
		dec, err := base64.StdEncoding.DecodeString(payload.Attachment)
		if err != nil {
			log.Error("Error when decoding the attachment", nlogger.Error(err))
			return nil, nil, ncore.TraceError(err)
		}

		// Create file attachment
		name := fmt.Sprintf("tmp/tmp.%v", payload.MimeType)
		f, err := os.Create(name)
		if err != nil {
			log.Error("Error occurred when create temp file", nlogger.Error(err))
			return nil, nil, ncore.TraceError(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			log.Error("Error occurred when write the temp file", nlogger.Error(err))
			return nil, nil, ncore.TraceError(err)
		}

		// Ensure that all the contents you've written are actually stored.
		if err := f.Sync(); err != nil {
			log.Error("Error occurred when sync the file", nlogger.Error(err))
			return nil, nil, ncore.TraceError(err)
		}

		// Attach file to message
		mailer.Attach(name)
	}

	return mailer, dialer, nil
}
