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
	log    nlogger.Logger
}

func NewServiceContext(ctx context.Context, config *contract.Config) contract.ServiceContext {
	return &ServiceContext{
		config: config,
		ctx:    ctx,
		log:    nlogger.Get().NewChild(logger.Context(ctx)),
	}
}
