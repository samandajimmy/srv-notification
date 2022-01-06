package contract

import (
	"context"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type ServiceContextConstructor = func(ctx context.Context, config *Config, source Source) ServiceContext

type ServiceContext interface {
	SendEmail(payload dto.SendEmail) error
	SendPushNotificationByTarget(payload dto.SendPushNotification) error
}

func NewService(core *ncore.Core, config *Config, source Source, fn ServiceContextConstructor) (*Service, error) {
	return &Service{
		Core:          core,
		config:        config,
		repo:          source,
		constructorFn: fn,
	}, nil
}

type Service struct {
	*ncore.Core
	config        *Config
	repo          Source
	constructorFn ServiceContextConstructor
}

func (s *Service) WithContext(ctx context.Context) ServiceContext {
	return s.constructorFn(ctx, s.config, s.repo)
}
