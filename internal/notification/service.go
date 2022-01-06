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
	repo   *RepositoryContext
	log    nlogger.Logger
}

func NewServiceContext(ctx context.Context, config *contract.Config) contract.ServiceContext {
	repo, err := NewRepository(&config.DataSources)
	if err != nil {
		return nil
	}

	return &ServiceContext{
		config: config,
		ctx:    ctx,
		repo:   repo.WithContext(ctx),
		log:    nlogger.Get().NewChild(logger.Context(ctx)),
	}
}
