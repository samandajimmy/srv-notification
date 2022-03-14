package notification

import (
	"context"
	"github.com/nbs-go/nlogger/v2"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type ServiceContext struct {
	config    *contract.Config
	ctx       context.Context
	repo      *RepositoryContext
	log       nlogger.Logger
	responses *ncore.ResponseMap
}

func NewServiceContext(ctx context.Context, config *contract.Config, core *ncore.Core) contract.ServiceContext {
	repo, err := NewRepository(config)
	if err != nil {
		log.Errorf("error create repository: %v", err)
		return nil
	}

	return &ServiceContext{
		config:    config,
		ctx:       ctx,
		repo:      repo.WithContext(ctx),
		log:       nlogger.Get().NewChild(logOption.Context(ctx)),
		responses: core.Responses,
	}
}

func (s *ServiceContext) Close() {
	// Close database connection to free pool
	err := s.repo.conn.Close()
	if err != nil {
		s.log.Error("Failed to close connection", logOption.Error(err))
	}
}
