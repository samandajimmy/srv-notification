package pds_svc

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	_ "github.com/lib/pq"
	"net/http"
	"path"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/repository"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/service"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nlogger"
)

var log = nlogger.Get()

type API struct {
	*contract.PdsApp
	dataSources repository.DataSourceMap
	PubSub      *gochannel.GoChannel
}

func NewAPI(core *ncore.Core, config contract.Config) API {
	// Init pubsub
	subLogger := watermill.NewStdLogger(true, true)
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, subLogger)

	return API{
		PdsApp: &contract.PdsApp{
			Core:         core,
			Config:       config,
			Repositories: contract.RepositoryMap{},
			Services: contract.ServiceMap{
				Auth:         new(service.Auth),
				Email:        new(service.Email),
				Notification: new(service.Notification),
			},
		},
		dataSources: repository.NewDataSourceMap(),
		PubSub:      pubSub,
	}
}

func (a *API) Boot() error {
	// Set value default configs
	a.Config.LoadFromEnv()

	// Init data sources
	err := a.dataSources.Init(a.Config.DataSources)
	if err != nil {
		return err
	}

	// Init repositories
	err = ncore.InitStruct(&a.Repositories, a.initRepository)
	if err != nil {
		return err
	}

	// Init services
	err = a.InitService()
	if err != nil {
		return err
	}

	return nil
}

func (a *API) InitRouter() http.Handler {
	// Init router
	router := nhttp.NewRouter(nhttp.RouterOptions{
		LogRequest: true,
		Debug:      a.Config.Server.Debug,
		TrustProxy: a.Config.Server.TrustProxy,
	})

	// Init handlers
	handlers := initHandler(a)

	// Enable cors
	if a.Config.CORS.Enabled {
		log.Debug("CORS Enabled")
		router.Use(a.Config.CORS.NewMiddleware())
	}

	// Set-up Routes
	setUpRoute(router, handlers)

	// Set-up Static
	staticPath := path.Join(a.WorkDir, "/web/static")
	staticDir := http.Dir(staticPath)
	staticServer := http.FileServer(staticDir)
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", staticServer))

	return router
}

func (a *API) InitSubscriber() {
	setUpSubscriber(a.PubSub, a.Services)
}

func (a *API) setUpStaticRoute(r *nhttp.Router) {
	staticServer := http.FileServer(http.Dir(path.Join(a.WorkDir, "/web/static")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticServer))
}

func (a *API) initRepository(name string, i interface{}) error {
	// Check interface
	r, ok := i.(repository.Initializer)
	if !ok {
		return fmt.Errorf("repository '%s' does not implement repository.Initializer interface", name)
	}

	// Init repository
	err := r.Init(a.dataSources, a.Repositories)
	if err != nil {
		return err
	}

	log.Debugf("Repositories.%s has been initialized", name)

	return nil
}
