package notification

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"gopkg.in/gomail.v2"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"

	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

const tmpAttachmentDir = ".tmp/email-attachments"

func (s *ServiceContext) SendEmail(payload *dto.SendEmail) error {
	var config model.SMTPConfig
	err := s.loadClientConfig(payload.ApplicationId, constant.SMTP, &config)
	if err != nil {
		return errx.Trace(err)
	}

	// Init mail client
	mailClient := s.newMailClient(config)

	// Set message
	message, tmpAttachment, err := s.composeEmail(payload)
	if err != nil {
		s.log.Error("Error when trying to compose email", logOption.Error(err))
		return errx.Trace(err)
	}

	// Send email
	err = mailClient.DialAndSend(message)
	if err != nil {
		s.log.Error("Failed dial and send email", logOption.Error(err))
		return errx.Trace(err)
	}

	// Clear attachment
	if tmpAttachment != "" {
		s.deleteEmailAttachmentTempFile(tmpAttachment)
	}

	s.log.Infof("Email has been sent successfully")
	return nil
}

func (s *ServiceContext) newMailClient(config model.SMTPConfig) *gomail.Dialer {
	return gomail.NewDialer(config.Host, nval.ParseIntFallback(config.Port, 465), config.Username, config.Password)
}

func (s *ServiceContext) composeEmail(payload *dto.SendEmail) (*gomail.Message, string, error) {
	// Set message
	msg := gomail.NewMessage()

	// Set sender format
	from := payload.From.Email
	if payload.From.Name != "" {
		from = fmt.Sprintf("%s <%s>", payload.From.Name, payload.From.Email)
	}

	// Set headers
	msg.SetHeader("From", from)
	msg.SetHeader("To", payload.To)
	msg.SetHeader("Subject", payload.Subject)
	msg.SetBody("text/html", payload.Message)

	// If no attachment, then return
	if payload.Attachment == "" {
		return msg, "", nil
	}

	// Write file attachment to temporary dir
	attachmentFileName, err := s.writeEmailAttachmentTempFile(payload.MimeType, payload.Attachment)
	if err != nil {
		return nil, "", errx.Trace(err)
	}

	// Attach file
	msg.Attach(attachmentFileName)

	// Return
	return msg, attachmentFileName, nil
}

func (s *ServiceContext) writeEmailAttachmentTempFile(fileExt string, b64FileContent string) (string, error) {
	fileContent, err := base64.StdEncoding.DecodeString(b64FileContent)
	if err != nil {
		s.log.Error("Error when decoding the attachment", logOption.Error(err))
		return "", errx.Trace(err)
	}

	// Generate temporary id
	tmpId, err := uuid.NewUUID()
	if err != nil {
		return "", errx.Trace(err)
	}

	// Compose file name
	fileName := fmt.Sprintf("%s/%s.%s", tmpAttachmentDir, tmpId, fileExt)

	// Create file
	f, err := os.Create(fileName)
	if err != nil {
		s.log.Error("Error occurred when create temp file", logOption.Error(err))
		return "", errx.Trace(err)
	}
	defer s.closeFile(f)

	// Write file content
	if _, err = f.Write(fileContent); err != nil {
		s.log.Error("Error occurred when write the temp file", logOption.Error(err))
		return "", errx.Trace(err)
	}

	// Ensure that all the contents you've written are actually stored.
	if err = f.Sync(); err != nil {
		s.log.Error("Error occurred when sync the file", logOption.Error(err))
		return "", errx.Trace(err)
	}

	return fileName, nil
}

func (s *ServiceContext) closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		s.log.Error("failed to close file", logOption.Error(err))
	}
}

func (s *ServiceContext) deleteEmailAttachmentTempFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		s.log.Error("Failed remove file after email sent!", logOption.Error(err))
	}
}
