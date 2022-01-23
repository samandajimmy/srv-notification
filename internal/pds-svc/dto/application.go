package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type Application struct {
	RequestId string   `json:"requestId"`
	Name      string   `json:"name"`
	Subject   *Subject `json:"-"`
}

func (d Application) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
	)
}

type GetApplication struct {
	RequestId string `json:"requestId"`
	XID       string `json:"xid"`
}

type AuthApplicationResponse struct {
	ID     int64  `json:"id"`
	XID    string `json:"XID"`
	Name   string `json:"name"`
	ApiKey string `json:"apiKey"`
}

type ApplicationResponse struct {
	XID  string `json:"xid"`
	Name string `json:"name"`
	ItemMetadataResponse
}

type ApplicationFindOptions struct {
	FindOptions
}

type ListApplicationResponse struct {
	Items    []*ApplicationResponse `json:"items"`
	Metadata ListMetadata           `json:"metadata"`
}

type ApplicationUpdateOptions struct {
	RequestId string          `json:"requestId"`
	XID       string          `json:"-"`
	Subject   *Subject        `json:"-"`
	Data      *Application    `json:"data"`
	Changelog map[string]bool `json:"changelog"`
	Version   int64           `json:"version"`
}

func (d ApplicationUpdateOptions) Validate() error {
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

	if err = Data(d.Data); err != nil {
		return ncore.TraceError(err)
	}

	return err
}

func Data(p *Application) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required, is.ASCII),
	)

	if err != nil {
		return nhttp.BadRequestError.Wrap(err)
	}

	return err
}
