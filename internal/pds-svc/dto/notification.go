package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SendPushNotification struct {
	RequestId     string            `json:"requestId"`
	Title         string            `json:"title"`
	Body          string            `json:"body"`
	ImageURL      string            `json:"imageUrl"`
	Token         string            `json:"token"`
	ApplicationId int64             `json:"applicationId"`
	Data          map[string]string `json:"data"`
}

func (d SendPushNotification) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ApplicationId, validation.Required),
		validation.Field(&d.Title, validation.Required),
		validation.Field(&d.Body, validation.Required),
		validation.Field(&d.Token, validation.Required),
	)
}

type AdditionalButton struct {
	ButtonLabel   string `json:"buttonLabel"`
	TransactionId string `json:"transactionId"`
	ScreenName    string `json:"screenName"`
}

type FCMOption struct {
	ImageUrl         string            `json:"imageUrl"`
	Token            string            `json:"token"`
	AdditionalButton AdditionalButton  `json:"additionalButton"`
	Data             map[string]string `json:"data"`
}

type SMTPOption struct {
	Subject    string     `json:"subject"`
	From       FromFormat `json:"from"`
	To         string     `json:"to"`
	Attachment string     `json:"attachment"`
	MimeType   string     `json:"mimeType"`
}

type NotificationOptionVO struct {
	FCM  *FCMOption  `json:"fcm"`
	SMTP *SMTPOption `json:"smtp"`
}

type SendNotificationOptionsRequest struct {
	RequestId      string                   `json:"-"`
	Auth           *AuthApplicationResponse `json:"auth"`
	Options        NotificationOptionVO     `json:"options"`
	UserId         int64                    `json:"userId"`
	Title          string                   `json:"title"`
	Content        string                   `json:"content"`
	ContentEncoded string                   `json:"contentEncoded"`
	ContentShort   string                   `json:"contentShort"`
}

func (d SendNotificationOptionsRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Options, validation.Required),
		validation.Field(&d.Title, validation.Required),
		validation.Field(&d.Content, validation.Required),
		validation.Field(&d.ContentEncoded, validation.Required, is.Base64),
		validation.Field(&d.ContentShort, validation.Required),
		validation.Field(&d.UserId, validation.Required),
	)
}
