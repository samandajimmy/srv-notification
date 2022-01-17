package notification

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"strings"
)

type ServiceContext struct {
	config    *contract.Config
	ctx       context.Context
	repo      *RepositoryContext
	log       nlogger.Logger
	responses *ncore.ResponseMap
}

func NewServiceContext(ctx context.Context, config *contract.Config, core *ncore.Core) contract.ServiceContext {
	repo, err := NewRepository(&config.DataSources)
	if err != nil {
		log.Errorf("error create repository: %v", err)
		return nil
	}

	return &ServiceContext{
		config:    config,
		ctx:       ctx,
		repo:      repo.WithContext(ctx),
		log:       nlogger.Get().NewChild(logger.Context(ctx)),
		responses: core.Responses,
	}
}

func (s *ServiceContext) GetOrderBy(sortBy string, sortDirection string, rules []string) (string, string) {
	if nval.InArrayString(sortBy, rules) {
		// Normalize direction
		sortDirection = strings.ToUpper(sortDirection)
		if sd := sortDirection; sd != `ASC` && sd != `DESC` {
			sortDirection = `ASC`
		}
		return sortBy, sortDirection
	}

	return `createdAt`, `DESC`
}
