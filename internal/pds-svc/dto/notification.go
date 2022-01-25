package dto

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SendPushNotificationRequest struct {
	RequestId string                   `json:"requestId"`
	Auth      *AuthApplicationResponse `json:"auth"`
	UserId    int64                    `json:"userId"`
	Title     string                   `json:"title"`
	Body      string                   `json:"body"`
	ImageUrl  string                   `json:"imageUrl"`
	Token     string                   `json:"token"`
	Metadata  json.RawMessage          `json:"metadata"`
	Data      map[string]string        `json:"data"`
}

func (d SendPushNotificationRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.Required),
		validation.Field(&d.Title, validation.Required),
		validation.Field(&d.Body, validation.Required),
		validation.Field(&d.Token, validation.Required),
	)
}

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
	UserId   int64             `json:"userId"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	ImageUrl string            `json:"imageUrl"`
	Token    string            `json:"token"`
	Metadata json.RawMessage   `json:"metadata"`
	Data     map[string]string `json:"data"`
}

type SMTPOption struct {
	UserId     int64      `json:"userId"`
	Subject    string     `json:"subject"`
	Message    string     `json:"message"`
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
	RequestId string                   `json:"-"`
	Auth      *AuthApplicationResponse `json:"auth"`
	Options   NotificationOptionVO     `json:"options"`
	UserId    int64                    `json:"userId"`
}

func (d SendNotificationOptionsRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Options, validation.Required),
		validation.Field(&d.UserId, validation.Required),
	)
}
