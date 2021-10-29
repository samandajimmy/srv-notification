package repository

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
)

type Initializer interface {
	ncore.InitializeChecker
	Init(dataSources DataSourceMap, repositories contract.RepositoryMap) error
}
