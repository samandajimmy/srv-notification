package pds_svc

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/handler"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"time"
)

type HandlerMap struct {
	Common       *handler.Common
	Email        *handler.Email
	Notification *handler.Notification
	Application  *handler.Application
}

func InitHandler(manifest *ncore.Manifest, svc *contract.Service, pubSub message.Publisher) *HandlerMap {
	return &HandlerMap{
		Common:       handler.NewCommon(time.Now(), manifest.AppVersion, manifest.GetStringMetadata(constant.BuildHashKey)),
		Email:        handler.NewEmail(pubSub),
		Notification: handler.NewNotification(pubSub),
		Application:  handler.NewApplication(svc),
	}
}
