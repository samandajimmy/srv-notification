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

	// Application
	router.Handle(http.MethodPost, "/applications",
		router.HandleFunc(handlers.Common.ValidateClient),
		router.HandleFunc(handlers.Application.PostApplication))
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
