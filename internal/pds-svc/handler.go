package pds_svc

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/handler"
	"time"
)

type HandlerMap struct {
	Common       *handler.Common
	Email        *handler.Email
	Notification *handler.Notification
}

func initHandler(app *API) *HandlerMap {

	return &HandlerMap{
		Common:       handler.NewCommon(time.Now(), app.Manifest.AppVersion, app.Manifest.GetStringMetadata(constant.BuildHashKey)),
		Email:        handler.NewEmail(app.Services.Email, app.PubSub),
		Notification: handler.NewNotification(app.Services.Notification, app.PubSub),
	}
}
