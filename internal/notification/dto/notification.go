package dto

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
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
	RequestId    string                      `json:"requestId"`
	Auth         *AuthApplicationResponse    `json:"auth"`
	Options      NotificationOptionVO        `json:"options"`
	UserId       int64                       `json:"userId"`
	Notification *DetailNotificationResponse `json:"notification"`
	*Subject
}

func (d SendNotificationOptionsRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Options, validation.Required),
		validation.Field(&d.UserId, validation.Required),
	)
}

type GetNotification struct {
	RequestId   string                   `json:"requestId"`
	ID          string                   `json:"id"`
	Application *AuthApplicationResponse `json:"application"`
}

func (d GetNotification) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ID, validation.Required, is.UUID),
		validation.Field(&d.Application, validation.Required),
	)
}

type GetCountNotification struct {
	RequestId   string                   `json:"requestId"`
	Application *AuthApplicationResponse `json:"applicationId"`
	UserRefId   int64                    `json:"UserRefId"`
}

func (d GetCountNotification) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Application, validation.Required),
		validation.Field(&d.UserRefId, validation.Required),
	)
}

type UpdateIsReadNotification struct {
	RequestId   string                   `json:"requestId"`
	Application *AuthApplicationResponse `json:"applicationId"`
	ID          string                   `json:"id"`
	Subject     *Subject                 `json:"-"`
}

func (d UpdateIsReadNotification) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ID, validation.Required, is.UUID),
		validation.Field(&d.Application, validation.Required),
	)
}

type DetailNotificationResponse struct {
	Id            uuid.UUID       `json:"id"`
	ApplicationId int64           `json:"applicationId"`
	UserRefId     int64           `json:"userRefId"`
	IsRead        bool            `json:"isRead"`
	ReadAt        int64           `json:"readAt"`
	Options       json.RawMessage `json:"options"`
	*BaseField
}

type DetailCountNotificationResponse struct {
	Count int64 `json:"count"`
}

type ListNotificationResponse struct {
	Items    []*DetailNotificationResponse `json:"items"`
	Metadata *ListMetadata                 `json:"metadata"`
}

type WebhookOptions struct {
	WebhookURL         string
	NotificationType   constant.NotificationType
	NotificationStatus constant.NotificationStatus
	Notification       *DetailNotificationResponse
	Payload            interface{}
}
