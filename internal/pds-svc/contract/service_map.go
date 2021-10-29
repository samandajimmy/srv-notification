package contract

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
)

type ServiceInitializer interface {
	ncore.InitializeChecker
	Init(app *PdsApp) error
}

type ServiceMap struct {
	Auth       AuthService
}
