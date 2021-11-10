package pds_svc

import (
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Send Email
	router.Handle(http.MethodPost, "/send-email", router.HandleFunc(handlers.Email.PostEmail))

	// Send Notification
	router.Handle(http.MethodPost, "/push-notification", router.HandleFunc(handlers.Notification.PostNotification))
}
