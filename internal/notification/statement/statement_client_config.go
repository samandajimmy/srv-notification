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

// Define schema

var ClientConfigSchema = schema.New(schema.FromModelRef(model.ClientConfig{}))

type ClientConfig struct {
	FindByKey                *sqlx.Stmt
	FindDefaultByKey         *sqlx.Stmt
	FindByXID                *sqlx.Stmt
	Insert                   *sqlx.NamedStmt
	UpdateByID               *sqlx.Stmt
	DeleteByID               *sqlx.Stmt
	FindJoinApplicationByXID *sqlx.Stmt
	IsExistsByKey            *sqlx.Stmt
	DeleteByApplicationID    *sqlx.Stmt
}

func NewClientConfig(db *nsql.Database) *ClientConfig {
	// Init query Schema Builder
	bs := query.Schema(ClientConfigSchema)

	// Init query
	findJoinApplicationByXID := query.Select(
		query.Column("*", option.Schema(ClientConfigSchema)),
		query.Column("*", option.Schema(ApplicationSchema))).
		From(ClientConfigSchema, option.As("cc")).
		Join(ApplicationSchema, query.Equal(query.Column("applicationId"), query.On("id")),
			option.As("a"), option.JoinMethod(op.InnerJoin)).
		Where(query.Equal(query.Column("xid"))).
		Limit(1).
		Build()

	findByKey := query.Select(query.Column("*")).
		Where(
			query.And(
				query.Equal(query.Column("key")),
				query.Equal(query.Column("applicationId")),
			),
		).
		From(ClientConfigSchema).
		Limit(1).
		Build()

	findByXID := query.Select(query.Column("*")).
		From(ClientConfigSchema).
		Where(query.Equal(query.Column("xid"))).
		Limit(1).
		Build()

	update := query.Update(ClientConfigSchema, "updatedAt", "modifiedBy", "version", "value").
		Where(query.And(
			query.Equal(query.Column("id")),
			query.Equal(query.Column("version")),
		)).
		Build(option.VariableFormat(op.BindVar))

	isExistsByKey := bs.IsExists(query.And(
		query.Equal(query.Column("applicationId")),
		query.Equal(query.Column("key")),
	))

	findDefaultByKey := query.Select(
		query.Column("*", option.Schema(ClientConfigSchema))).
		From(ClientConfigSchema, option.As("cc")).
		Join(ApplicationSchema, query.And(
			query.Equal(query.Column("applicationId"), query.On("id")),
			query.Equal(query.Column("xid", option.Schema(ApplicationSchema), query.BindVar())),
		), option.As("a"), option.JoinMethod(op.InnerJoin)).
		Where(query.Equal(query.Column("key"))).
		Limit(1).
		Build()

	deleteByAppId := query.Delete(ClientConfigSchema).
		Where(query.Equal(query.Column("applicationId"))).
		Build()

	return &ClientConfig{
		FindByKey:                db.PrepareRebind(findByKey),
		FindDefaultByKey:         db.PrepareRebind(findDefaultByKey),
		FindByXID:                db.PrepareRebind(findByXID),
		Insert:                   db.PrepareNamed(bs.Insert()),
		UpdateByID:               db.PrepareRebind(update),
		DeleteByID:               db.PrepareRebind(bs.Delete()),
		FindJoinApplicationByXID: db.PrepareRebind(findJoinApplicationByXID),
		IsExistsByKey:            db.PrepareRebind(isExistsByKey),
		DeleteByApplicationID:    db.PrepareRebind(deleteByAppId),
	}
}
