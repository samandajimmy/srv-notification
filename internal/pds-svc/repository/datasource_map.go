package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type DataSourceMap struct {
	Postgres *nsql.DB
}

func (a *DataSourceMap) Init(config contract.DataSourcesConfig) error {
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
