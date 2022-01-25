package pds_svc

import (
	"github.com/nbs-go/nlogger"
	"net/http"
	"path"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}

func setUpRoute(router *nhttp.Router, handlers *HandlerMap) {
	// Common
	router.Handle(http.MethodGet, "/", router.HandleFunc(handlers.Common.GetAPIStatus))

	// Send Email
	router.Handle(http.MethodPost, "/send-email", router.HandleFunc(handlers.Email.PostEmail))

	// Send Notification
	router.Handle(http.MethodPost, "/push-notification", router.HandleFunc(handlers.Notification.PostNotification))
	router.Handle(http.MethodPost, "/notifications",
		router.HandleFunc(handlers.Middlewares.AuthApp),
		router.HandleFunc(handlers.Notification.PostCreateNotification))

	// Application
	router.Handle(http.MethodPost, "/applications",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.Application.PostCreateApplication))
	router.Handle(http.MethodGet, "/applications", router.HandleFunc(handlers.Application.GetFindApplication))
	router.Handle(http.MethodGet, "/applications/{xid}", router.HandleFunc(handlers.Application.GetDetailApplication))
	router.Handle(http.MethodPut, "/applications/{xid}",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.Application.PutUpdateApplication))
	router.Handle(http.MethodDelete, "/applications/{xid}", router.HandleFunc(handlers.Application.DeleteApplication))

	// Client Config
	router.Handle(http.MethodGet, "/client-configs",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.ClientConfig.SearchClientConfig))
	router.Handle(http.MethodGet, "/client-configs/{xid}", router.HandleFunc(handlers.ClientConfig.DetailClientConfig))
	router.Handle(http.MethodPost, "/client-configs", router.HandleFunc(handlers.Common.ValidateClient), router.HandleFunc(handlers.ClientConfig.CreateClientConfig))
	router.Handle(http.MethodPut, "/client-configs/{xid}",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.ClientConfig.UpdateClientConfig))
	router.Handle(http.MethodDelete, "/client-configs/{xid}",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.ClientConfig.DeleteClientConfig))

	// Notification
	router.Handle(http.MethodGet, "/notifications/{id}", router.HandleFunc(handlers.Notification.GetDetailNotification))
	router.Handle(http.MethodDelete, "/notifications/{id}", router.HandleFunc(handlers.Notification.DeleteNotification))
}

func InitRouter(workDir string, config *contract.Config, handlers *HandlerMap) http.Handler {
	// Init router
	router := nhttp.NewRouter(nhttp.RouterOptions{
		LogRequest: true,
		Debug:      config.Server.Debug,
		TrustProxy: config.Server.TrustProxy,
	})

	// Enable cors
	if config.CORS.Enabled {
		log.Debug("CORS Enabled")
		router.Use(config.CORS.NewMiddleware())
	}

	// Set-up Routes
	setUpRoute(router, handlers)

	// Set-up Static
	staticPath := path.Join(workDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}
