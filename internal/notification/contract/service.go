package contract

import (
	"context"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type ServiceContextConstructor = func(ctx context.Context, config *Config, core *ncore.Core) ServiceContext

type ServiceContext interface {
	Close()

	AuthApplication(username string, password string) (*dto.AuthApplicationResponse, error)

	SendEmail(payload *dto.SendEmail) error
	SendPushNotificationByTarget(payload *dto.SendPushNotification) error

	CreateApplication(payload *dto.Application) (*dto.ApplicationResponse, error)
	GetDetailApplication(payload *dto.GetApplication) (*dto.ApplicationResponse, error)
	DeleteApplication(payload *dto.GetApplication) error
	ListApplication(options *dto.ListPayload) (*dto.ListApplicationResponse, error)
	UpdateApplication(payload *dto.ApplicationUpdateOptions) (*dto.ApplicationResponse, error)

	CreateClientConfig(payload *dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error)
	GetDetailClientConfig(payload *dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error)
	ListClientConfig(params *dto.ListPayload) (*dto.ClientConfigListResponse, error)
	DeleteClientConfig(payload *dto.GetClientConfig) error
	UpdateClientConfig(payload *dto.ClientConfigUpdateOptions) (*dto.ClientConfigItemResponse, error)

	CreateNotification(payload *dto.SendNotificationOptionsRequest) (*dto.DetailNotificationResponse, error)
	GetDetailNotification(payload *dto.GetNotification) (*dto.DetailNotificationResponse, error)
	CountNotification(payload *dto.GetCountNotification) (*dto.DetailCountNotificationResponse, error)
	DeleteNotification(payload *dto.GetNotification) error
	ListNotification(options *dto.ListPayload) (*dto.ListNotificationResponse, error)
	UpdateIsRead(payload *dto.UpdateIsReadNotification) (*dto.DetailNotificationResponse, error)
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
	return s.constructorFn(ctx, s.config, s.Core)
}
