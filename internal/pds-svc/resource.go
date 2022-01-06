package pds_svc

import (
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/repository"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type Resource struct {
	*contract.Source
	Config     contract.Config
	DataSource repository.DataSourceMap
}

func NewResource(config contract.Config) Resource {
	return Resource{
		&contract.Source{
			Repositories: contract.RepositoryMap{
				Configure: new(repository.ClientConfig),
			},
		},
		config,
		repository.NewDataSourceMap(),
	}
}

func (a *Resource) BootResource() error {
	// Init data sources
	err := a.DataSource.Init(a.Config.DataSources)
	if err != nil {
		return err
	}

	// Init repositories
	err = ncore.InitStruct(&a.Repositories, a.initRepository)
	if err != nil {
		return err
	}

	return nil
}

func (a *Resource) initRepository(name string, i interface{}) error {
	// Check interface
	r, ok := i.(repository.Initializer)
	if !ok {
		return fmt.Errorf("repository '%s' does not implement repository.Initializer interface", name)
	}

	// Init repository
	err := r.Init(a.DataSource, a.Repositories)
	if err != nil {
		return err
	}

	log.Debugf("Repositories.%s has been initialized", name)

	return nil
}
