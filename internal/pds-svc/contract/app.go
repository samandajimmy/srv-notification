package contract

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nlogger"
	"fmt"
)

var log = nlogger.Get()

type PdsApp struct {
	*ncore.Core
	Config       Config
	Repositories RepositoryMap
	Services     ServiceMap
}

func (a *PdsApp) InitService() error {
	// Init services
	err := ncore.InitStruct(&a.Services, a.initService)
	return err
}

func (a *PdsApp) initService(name string, i interface{}) error {
	// Check interface
	r, ok := i.(ServiceInitializer)
	if !ok {
		return fmt.Errorf("service '%s' does not implement ServiceInitializer interface", name)
	}

	// Init service
	err := r.Init(a)
	if err != nil {
		return err
	}

	log.Debugf("Services.%s has been initialized", name)

	return nil
}
