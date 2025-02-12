package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SendEmail struct {
	ApplicationId int64      `json:"applicationId"`
	RequestId     string     `json:"requestId"`
	Subject       string     `json:"subject"`
	From          FromFormat `json:"from"`
	To            string     `json:"to"`
	Message       string     `json:"message"`
	Attachment    string     `json:"attachment"`
	MimeType      string     `json:"mimeType"`
}

func (d SendEmail) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ApplicationId, validation.Required),
		validation.Field(&d.Subject, validation.Required),
		validation.Field(&d.From),
		validation.Field(&d.To, validation.Required, is.EmailFormat),
		validation.Field(&d.Message, validation.Required),
	)
}

type FromFormat struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (d FromFormat) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Email, validation.Required, is.EmailFormat),
	)
}
