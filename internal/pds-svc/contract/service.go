package contract

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
)

type AuthService interface {
	ValidateClient(payload dto.ClientCredential) error
}

type EmailService interface {
	SendEmail(payload dto.SendEmail) error
}

type NotificationService interface {
	SendNotificationByToken(payload dto.NotificationCreate) error
}
