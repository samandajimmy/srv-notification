package main

import (
	// Inject logger before loading other packages
	_ "repo.pegadaian.co.id/ms-pds/srv-notification/cmd/notification/logger"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbs-go/nlogger/v2"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"net/http"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"time"
)

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}

func main() {
	// Capture started at
	startedAt := time.Now()

	// Init config from env
	config := new(contract.Config)
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	// Boot core
	bootOptions := handleCmdFlags()
	core := ncore.Boot(bootOptions.Core)

	// Check if migration option is set
	err = bootMigration(core.WorkDir, config)
	if err != nil {
		panic(err)
	}

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

	// Start server
	log.Infof("Starting %s...", core.Manifest.AppName)
	log.Infof("NodeId = %s, Environment = %s", core.NodeId, core.GetEnvironmentString())
	log.Debugf("Boot Time: %s", time.Since(startedAt))

	err = http.ListenAndServe(nhttp.ListenPort(config.Port), router)
	if err != nil {
		log.Fatal("failed to serve", logOption.Error(err))
		os.Exit(2)
	}
}
