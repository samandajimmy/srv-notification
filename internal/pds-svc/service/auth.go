package service

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
)

type Auth struct {
	clientSecret string
	clientID     string
	responses    *ncore.ResponseMap
}

func (s *Auth) ValidateClient(payload dto.ClientCredential) error {
	if payload.ClientID != s.clientID || payload.ClientSecret != s.clientSecret {
		return s.responses.GetError("E_AUTH_1")
	}

	return nil
}

func (s *Auth) HasInitialized() bool {
	return s.clientID != "" && s.clientSecret != ""
}

func (s *Auth) Init(app *contract.PdsApp) error {
	cfg := app.Config.Client
	s.clientID = cfg.ClientID
	s.clientSecret = cfg.ClientSecret
	s.responses = app.Responses
	return nil
}
