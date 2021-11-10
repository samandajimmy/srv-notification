package contract

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

type ServiceInitializer interface {
	ncore.InitializeChecker
	Init(app *PdsApp) error
}

type ServiceMap struct {
	Auth         AuthService
	Email        EmailService
	Notification NotificationService
}
