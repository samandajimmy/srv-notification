package dto

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ClientConfig struct {
	RequestId     string            `json:"requestId"`
	XID           string            `json:"xid"`
	Key           string            `json:"key"`
	Value         map[string]string `json:"value"`
	ApplicationId int               `json:"applicationId"`
	Subject       *Subject          `json:"-"`
}

type GetClientConfig struct {
	XID       string `json:"xid"`
	RequestId string `json:"requestId"`
}

func (d ClientConfig) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.XID, validation.Required),
		validation.Field(&d.Key, validation.Required),
		validation.Field(&d.Value, validation.Required),
		validation.Field(&d.ApplicationId, validation.Required),
	)
}

type ClientConfigItemResponse struct {
	XID           string          `json:"xid"`
	Key           string          `json:"key"`
	Value         json.RawMessage `json:"value"`
	ApplicationId int             `json:"applicationId"`
	ItemMetadataResponse
}

type ClientConfigListResponse struct {
	ClientConfig []ClientConfigItemResponse `json:"rows"`
	Metadata     ListMetadata               `json:"metadata"`
}

type ClientConfigFindOptions struct {
	FindOptions
	Subject *Subject
}
