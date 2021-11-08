package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NotificationCreate struct {
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageURL string            `json:"imageUrl"`
	Token    string            `json:"token"`
	Data     map[string]string `json:"data"`
}

func (d NotificationCreate) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Title, validation.Required),
		validation.Field(&d.Body, validation.Required),
		validation.Field(&d.Token, validation.Required),
	)
}
