package nmail

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

type NucleoEmail struct {
	Client *gomail.Dialer
}

func NewNucleoEmail(credential contract.SMTPConfig) (*NucleoEmail, error) {

	// Set credential SMTP
	host := credential.Host
	port := credential.Port
	username := credential.Username
	password := credential.Password

	// Set dialer
	dialer := gomail.NewDialer(host, port, username, password)

	return &NucleoEmail{
		Client: dialer,
	}, nil
}

func (ne *NucleoEmail) ComposeEmail(payload dto.SendEmail) (*gomail.Message, error) {

	// Set message
	mailer := gomail.NewMessage()

	// Set sender format
	from := fmt.Sprintf("%v", payload.From.Email)
	if payload.From.Name != "" {
		from = fmt.Sprintf("%v <%v>", payload.From.Name, payload.From.Email)
	}

	mailer.SetHeader("From", from)
	mailer.SetHeader("To", payload.To)
	mailer.SetHeader("Subject", payload.Subject)
	mailer.SetBody("text/html", payload.Message)

	if payload.Attachment != "" {
		dec, err := base64.StdEncoding.DecodeString(payload.Attachment)
		if err != nil {
			log.Error("Error when decoding the attachment", nlogger.Error(err))
			return nil, ncore.TraceError(err)
		}

		// Create file attachment
		name := fmt.Sprintf("tmp/tmp.%v", payload.MimeType)
		f, err := os.Create(name)
		if err != nil {
			log.Error("Error occurred when create temp file", nlogger.Error(err))
			return nil, ncore.TraceError(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			log.Error("Error occurred when write the temp file", nlogger.Error(err))
			return nil, ncore.TraceError(err)
		}

		// Ensure that all the contents you've written are actually stored.
		if err := f.Sync(); err != nil {
			log.Error("Error occurred when sync the file", nlogger.Error(err))
			return nil, ncore.TraceError(err)
		}

		// Attach file to message
		mailer.Attach(name)
	}

	return mailer, nil
}
