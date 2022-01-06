package notification

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
)

type ServiceContext struct {
	config *contract.Config
	ctx    context.Context
	repo   contract.RepositoryMap
	log    nlogger.Logger
}

func NewServiceContext(ctx context.Context, config *contract.Config, source contract.Source) contract.ServiceContext {
	repo := source.Repositories
	return &ServiceContext{
		config: config,
		ctx:    ctx,
		repo:   repo,
		log:    nlogger.Get().NewChild(logger.Context(ctx)),
	}
}
