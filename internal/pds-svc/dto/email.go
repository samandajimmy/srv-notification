package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SendEmail struct {
	Subject    string `json:"subject"`
	From       string `json:"from"`
	To         string `json:"to"`
	Message    string `json:"message"`
	Attachment string `json:"attachment"`
	MimeType   string `json:"mimeType"`
}

func (d SendEmail) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Subject, validation.Required),
		validation.Field(&d.From, validation.Required),
		validation.Field(&d.To, validation.Required),
		validation.Field(&d.Message, validation.Required),
	)
}
