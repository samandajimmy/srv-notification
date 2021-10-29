package pds_svc

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/constant"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/handler"
	"time"
)

type HandlerMap struct {
	Common   *handler.Common
}

func initHandler(app *API) *HandlerMap {

	return &HandlerMap{
		Common:   handler.NewCommon(time.Now(), app.Manifest.AppVersion, app.Manifest.GetStringMetadata(constant.BuildHashKey)),
	}
}
