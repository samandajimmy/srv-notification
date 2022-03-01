package notification

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"

	_ "github.com/lib/pq"
)

func NewRepository(config *contract.Config) (*Repository, error) {
	// Parse options
	maxIdleConn := nval.ParseIntFallback(config.DatabaseMaxIdleConn, 10)
	maxOpenConn := nval.ParseIntFallback(config.DatabaseMaxOpenConn, 10)
	maxConnLifetime := nval.ParseIntFallback(config.DatabaseMaxConnLifetime, 1)

	// Init db
	db, err := nsql.NewDatabase(nsql.Config{
		Driver:          config.DatabaseDriver,
		Host:            config.DatabaseHost,
		Port:            config.DatabasePort,
		Username:        config.DatabaseUsername,
		Password:        config.DatabasePassword,
		Database:        config.DatabaseName,
		MaxIdleConn:     &maxIdleConn,
		MaxOpenConn:     &maxOpenConn,
		MaxConnLifetime: &maxConnLifetime,
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
