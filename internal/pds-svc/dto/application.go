package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

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

type ApplicationResponse struct {
	XID  string `json:"xid"`
	Name string `json:"name"`
	ItemMetadataResponse
}
