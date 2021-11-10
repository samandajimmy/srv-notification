package repository

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type Initializer interface {
	ncore.InitializeChecker
	Init(dataSources DataSourceMap, repositories contract.RepositoryMap) error
}
