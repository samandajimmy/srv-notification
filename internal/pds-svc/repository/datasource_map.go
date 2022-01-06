package repository

import (
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

var log = nlogger.Get()

type DataSourceMap struct {
	Postgres *nsql.DB
}

func (a *DataSourceMap) Init(config contract.DataSourcesConfig) error {
	// Skip if not initialized
	if a.Postgres == nil {
		log.Debug("Skipping db init")
		return nil
	}

	err := a.Postgres.Init(config.Postgres)
	if err != nil {
		return ncore.TraceError(err)
	}

	return nil
}

func NewDataSourceMap() DataSourceMap {
	return DataSourceMap{
		Postgres: new(nsql.DB),
	}
}
