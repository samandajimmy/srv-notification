package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/nbs-go/errx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type Application struct {
	RequestId  string   `json:"requestId"`
	Name       string   `json:"name"`
	ApiKey     string   `json:"apiKey"`
	WebhookURL string   `json:"webhookUrl"`
	Subject    *Subject `json:"-"`
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
	ID         int64  `json:"id"`
	XID        string `json:"xid"`
	Name       string `json:"name"`
	ApiKey     string `json:"apiKey"`
	WebhookURL string `json:"webhookUrl"`
}

type ApplicationResponse struct {
	XID        string `json:"xid"`
	Name       string `json:"name"`
	ApiKey     string `json:"apiKey"`
	WebhookURL string `json:"webhookUrl"`
	*BaseField
}

type ApplicationItem struct {
	XID  string `json:"xid"`
	Name string `json:"name"`
	*BaseField
}

type ListApplicationResponse struct {
	Items    []*ApplicationItem `json:"items"`
	Metadata *ListMetadata      `json:"metadata"`
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
		return errx.Trace(err)
	}

	return err
}

func Data(p *Application) error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required, is.ASCII),
		validation.Field(&p.WebhookURL, is.URL),
	)

	if err != nil {
		return nhttp.BadRequestError.Wrap(err)
	}

	return err
}
