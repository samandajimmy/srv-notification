package notification

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"

	_ "github.com/lib/pq"
)

var log = nlogger.Get()

func NewRepository(config *contract.DataSourcesConfig) (*Repository, error) {
	configPostgres := config.Postgres

	// Init db
	db, err := nsql.NewDatabase(nsql.Config{
		Driver:          configPostgres.Driver,
		Host:            configPostgres.Host,
		Port:            configPostgres.Port,
		Username:        configPostgres.Username,
		Password:        configPostgres.Password,
		Database:        configPostgres.Database,
		MaxIdleConn:     configPostgres.MaxIdleConn,
		MaxOpenConn:     configPostgres.MaxOpenConn,
		MaxConnLifetime: configPostgres.MaxConnLifetime,
	})
	if err != nil {
		log.Error("failed to initiate connection to db", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	// Init repo
	r := Repository{
		db: db,
	}

	return &r, nil
}

type Repository struct {
	db   *nsql.Database
	stmt *RepositoryStatement
}

type RepositoryContext struct {
	*RepositoryStatement
	ctx  context.Context
	conn *sqlx.Conn
}

func (r *Repository) WithContext(ctx context.Context) *RepositoryContext {
	// If db is not connected, then initialize connection
	isConnected, _ := r.db.IsConnected(ctx)
	if !isConnected {
		log.Debugf("initialize connection to database...")
		err := r.db.Init()
		if err != nil {
			log.Error("failed to initiate connection to db", nlogger.Error(err))
			panic(ncore.TraceError(err))
		}
	}

	// If statement has not been initiated, then init
	if r.stmt == nil {
		log.Debugf("initialize statement...")
		r.stmt = NewRepositoryStatement(r.db)
	}

	// Get connection
	conn, err := r.db.GetConnection(ctx)
	if err != nil {
		log.Error("failed to retrieve connection to db", nlogger.Error(err))
		panic(ncore.TraceError(err))
	}

	return &RepositoryContext{
		ctx:                 ctx,
		conn:                conn,
		RepositoryStatement: r.stmt,
	}
}

func (rc *RepositoryContext) GetOrderByQuery(sortBy string, sortDirection string) string {
	return fmt.Sprintf(`"%s" %s`, sortBy, sortDirection)
}
