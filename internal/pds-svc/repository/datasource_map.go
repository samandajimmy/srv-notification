package repository

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nsql"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
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
