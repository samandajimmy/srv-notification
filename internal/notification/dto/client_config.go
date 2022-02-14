package dto

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type ClientConfigUpdateOptions struct {
	RequestId string          `json:"requestId"`
	XID       string          `json:"-"`
	Subject   *Subject        `json:"-"`
	Data      *ClientConfig   `json:"data"`
	Changelog map[string]bool `json:"changelog"`
	Version   int64           `json:"version"`
}

type ClientConfigRequest struct {
	RequestId string   `json:"requestId"`
	XID       string   `json:"-"`
	Subject   *Subject `json:"-"`
	ClientConfig
}

type ClientConfig struct {
	ApplicationXid string            `json:"applicationXid"`
	Key            string            `json:"key"`
	Value          map[string]string `json:"value"`
}

type GetClientConfig struct {
	XID       string `json:"xid"`
	RequestId string `json:"requestId"`
}

func (d ClientConfig) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Key, validation.Required),
		validation.Field(&d.Value, validation.Required),
		validation.Field(&d.ApplicationXid, validation.Required),
	)
}

func (d ClientConfigUpdateOptions) Validate() error {
	err := validation.ValidateStruct(&d,
		validation.Field(&d.XID, validation.Required),
		validation.Field(&d.Subject, validation.Required),
		validation.Field(&d.Changelog, validation.Required),
		validation.Field(&d.Version, validation.Required, validation.Min(1)),
		validation.Field(&d.Data, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Wrap(err)
	}

	if err = ClientConfigDataValidate(d.Data); err != nil {
		return ncore.TraceError(err)
	}

	return err
}

type ClientConfigItemResponse struct {
	ApplicationXid string          `json:"applicationXid"`
	XID            string          `json:"xid"`
	Key            string          `json:"key"`
	Value          json.RawMessage `json:"value"`
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

func ClientConfigDataValidate(p *ClientConfig) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Key, validation.Required),
		validation.Field(&p.Value, validation.Required),
	)

	if err != nil {
		return nhttp.BadRequestError.Wrap(err)
	}

	return err
}
