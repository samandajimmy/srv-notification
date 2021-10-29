package pds_svc

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nhttp"
	"net/http"
)

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))
}
