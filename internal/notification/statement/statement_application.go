package statement

import (
	"github.com/jmoiron/sqlx"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

// Declare schema

var ApplicationSchema = schema.New(schema.FromModelRef(model.Application{}))

type Application struct {
	Insert     *sqlx.NamedStmt
	FindByXID  *sqlx.Stmt
	DeleteByID *sqlx.Stmt
	Update     *sqlx.Stmt
}

func NewApplication(db *nsql.Database) *Application {
	// Create schema builder
	sb := query.Schema(ApplicationSchema)

	// Create builder
	findByXID := query.Select(query.Column("*")).
		From(ApplicationSchema).
		Where(query.Equal(query.Column("xid"))).
		Build()

	// Update
	update := query.Update(ApplicationSchema, "updatedAt", "modifiedBy", "version", "name", "webhookUrl").
		Where(query.And(
			query.Equal(query.Column("id")),
			query.Equal(query.Column("version")),
		)).
		Build(option.VariableFormat(op.BindVar))

	// Update builder
	return &Application{
		Insert:     db.PrepareNamed(sb.Insert()),
		FindByXID:  db.PrepareRebind(findByXID),
		DeleteByID: db.PrepareRebind(sb.Delete()),
		Update:     db.PrepareRebind(update),
	}
}
