package notification

import (
	"github.com/nbs-go/nlogger/v2"
	"net/http"
	"path"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"strings"
)

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Application
	router.Handle(http.MethodPost, "/applications",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.Application.CreateApplication))
	router.Handle(http.MethodGet, "/applications",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.Application.ListApplication))
	router.Handle(http.MethodGet, "/applications/{xid}",
		router.HandleFunc(handlers.Application.GetDetailApplication))
	router.Handle(http.MethodPut, "/applications/{xid}",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.Application.UpdateApplication))
	router.Handle(http.MethodDelete, "/applications/{xid}",
		router.HandleFunc(handlers.Application.DeleteApplication))

	// Client Config
	router.Handle(http.MethodGet, "/client-configs",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.ClientConfig.ListClientConfig))
	router.Handle(http.MethodGet, "/client-configs/{xid}",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.ClientConfig.GetDetailClientConfig))
	router.Handle(http.MethodPost, "/client-configs",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.ClientConfig.CreateClientConfig))
	router.Handle(http.MethodPut, "/client-configs/{xid}",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.ClientConfig.UpdateClientConfig))
	router.Handle(http.MethodDelete, "/client-configs/{xid}",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.ClientConfig.DeleteClientConfig))

	// Notification
	router.Handle(http.MethodPost, "/notifications",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.PostCreateNotification))
	router.Handle(http.MethodGet, "/notifications/count",
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.CountNotification))
	router.Handle(http.MethodGet, "/notifications/{id}",
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.GetDetailNotification))
	router.Handle(http.MethodDelete, "/notifications/{id}",
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.DeleteNotification))
	router.Handle(http.MethodGet, "/notifications",
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.ListNotification))
	router.Handle(http.MethodPut, "/notifications/{id}/is-read",
		router.HandleFunc(handlers.Common.ParseSubject),
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.UpdateIsReadNotification))
}

func InitRouter(workDir string, config *contract.Config, handlers *HandlerMap) http.Handler {
	// Check debug
	debug := strings.ToLower(config.Debug) == "true"
	trustProxy := strings.ToLower(config.ServerTrustProxy) == "true"

	// Init router
	router := nhttp.NewRouter(nhttp.RouterOptions{
		LogRequest: true,
		Debug:      debug,
		TrustProxy: trustProxy,
	})

	// Set-up Routes
	setUpRoute(router, handlers)

	// Set-up Static
	staticPath := path.Join(workDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}
