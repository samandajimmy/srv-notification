package main

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/contract"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nlogger"
	"fmt"
	"net/http"
	"time"
)

var log = nlogger.Get()

func main() {
	// Boot
	bootStartedAt := time.Now()
	app := boot()
	log.Debugf("Boot Time: %s", time.Since(bootStartedAt))

	// Start app
	start(&app)
}

func boot() pds_svc.API {
	// Handle command
	bootOptions := handleCmdFlags()

	// Boot Core
	config := contract.Config{}
	core := ncore.Boot(&config, bootOptions.Core)

	// Boot App
	app := pds_svc.NewAPI(core, config)
	err := app.Boot()
	if err != nil {
		panic(err)
	}

	return app
}

func start(app *pds_svc.API) {
	// Set server config from env
	err := app.Config.Server.LoadFromEnv()
	if err != nil {
		panic(err)
	}

	// Get server config
	config := app.Config.Server

	// Init router
	router := app.InitRouter()

	log.Infof("%s HTTP Server is listening to port %d", AppSlug, config.ListenPort)
	log.Infof("%s HTTP Server Started. Base URL: %s", AppSlug, config.GetHttpBaseUrl())
	err = http.ListenAndServe(config.GetListenPort(), router)
	if err != nil {
		panic(fmt.Errorf("%s: failed to on listen.\n  > %w", AppSlug, err))
	}
}
