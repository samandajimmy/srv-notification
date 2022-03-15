package handler

import (
	"fmt"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

type Middlewares struct {
	svc *contract.Service
}

func NewMiddlewares(svc *contract.Service) *Middlewares {
	m := Middlewares{svc}
	return &m
}

func (h *Middlewares) AuthApp(rx *nhttp.Request) (*nhttp.Response, error) {
	// validate basic auth
	username, password, ok := rx.BasicAuth()
	if !ok {
		return nil, nhttp.ForbiddenError
	}

	// TODO: Refactor to Auth Middleware
	// Call service
	svc := h.svc.WithContext(rx.Context())
	defer svc.Close()

	app, err := svc.AuthApplication(username, password)
	if err != nil {
		log.Error("failed to call service.", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Set context
	rx.SetContextValue(constant.ApplicationKey, app)

	return nhttp.Continue(), nil
}

func getApplication(rx *nhttp.Request) (*dto.AuthApplicationResponse, error) {
	v := rx.GetContextValue(constant.ApplicationKey)

	app, ok := v.(*dto.AuthApplicationResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected app value in context. Type: %T", v)
	}

	return app, nil
}
