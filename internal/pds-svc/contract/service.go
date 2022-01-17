package contract

import (
	"context"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type ServiceContextConstructor = func(ctx context.Context, config *Config) ServiceContext

type ServiceContext interface {
	SendEmail(payload dto.SendEmail) error
	SendPushNotificationByTarget(payload dto.SendPushNotification) error
	CreateApplication(payload dto.Application) (*dto.ApplicationResponse, error)
}

func NewService(core *ncore.Core, config *Config, fn ServiceContextConstructor) (*Service, error) {
	return &Service{
		Core:          core,
		config:        config,
		constructorFn: fn,
	}, nil
}

type Service struct {
	*ncore.Core
	config        *Config
	constructorFn ServiceContextConstructor
}

func (s *Service) WithContext(ctx context.Context) ServiceContext {
	return s.constructorFn(ctx, s.config)
}
