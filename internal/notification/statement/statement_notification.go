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

var NotificationSchema = schema.New(schema.FromModelRef(model.Notification{}), schema.AutoIncrement(false))

type Notification struct {
	Insert     *sqlx.NamedStmt
	FindByID   *sqlx.Stmt
	Update     *sqlx.Stmt
	DeleteByID *sqlx.Stmt
}

func NewNotification(db *nsql.Database) *Notification {
	sb := query.Schema(NotificationSchema)

	update := query.Update(NotificationSchema, "updatedAt", "modifiedBy", "version", "isRead", "readAt").
		Where(query.And(
			query.Equal(query.Column("id")),
			query.Equal(query.Column("version")),
		)).
		Build(option.VariableFormat(op.BindVar))

	return &Notification{
		Insert:     db.PrepareNamed(sb.Insert()),
		FindByID:   db.PrepareRebind(sb.FindByPK()),
		Update:     db.PrepareRebind(update),
		DeleteByID: db.PrepareRebind(sb.Delete()),
	}
}
