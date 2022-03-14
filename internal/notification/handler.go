package notification

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/handler"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"time"
)

type HandlerMap struct {
	Middlewares  *handler.Middlewares
	Common       *handler.Common
	Notification *handler.Notification
	Application  *handler.Application
	ClientConfig *handler.ClientConfig
}

func InitHandler(manifest *ncore.Manifest, svc *contract.Service, pubSub message.Publisher) *HandlerMap {
	return &HandlerMap{
		Middlewares:  handler.NewMiddlewares(svc),
		Common:       handler.NewCommon(time.Now(), manifest.AppVersion, manifest.GetStringMetadata(constant.BuildHashKey)),
		Notification: handler.NewNotification(pubSub, svc),
		Application:  handler.NewApplication(svc),
		ClientConfig: handler.NewClientConfig(svc),
	}
}
