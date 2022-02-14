package main

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/nbs-go/nlogger"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	_ "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"time"
)

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}

func main() {
	// Capture started at
	startedAt := time.Now()

	// TODO: handle boot options
	bootOptions := handleCmdFlags()

	// Load config
	config := contract.LoadConfig()

	// Boot core
	core := ncore.Boot(bootOptions.Core)

	// Init service
	svc, err := contract.NewService(core, config, notification.NewServiceContext)
	if err != nil {
		panic(err)
	}

	// Init pubsub
	//subLogger := watermill.NewStdLogger(true, true)
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, nil)

	// Init subscriber
	notification.SetUpSubscriber(pubSub, svc)

	// Init handlers
	handlers := notification.InitHandler(&core.Manifest, svc, pubSub)

	log.Debugf("Boot Time: %s", time.Since(startedAt))

	// Init router
	router := notification.InitRouter(core.WorkDir, config, handlers)

	// Set server config from env
	err = config.Server.LoadFromEnv()
	if err != nil {
		panic(err)
	}
	serverConfig := config.Server

	// Start server
	log.Infof("%s HTTP Server is listening to port %d", AppSlug, serverConfig.ListenPort)
	log.Infof("%s HTTP Server Started. Base URL: %s", AppSlug, serverConfig.GetHttpBaseUrl())
	err = http.ListenAndServe(serverConfig.GetListenPort(), router)
	if err != nil {
		panic(fmt.Errorf("%s: failed on listen.\n  > %w", AppSlug, err))
	}
}
