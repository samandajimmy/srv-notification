package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type Application struct {
	Insert    *sqlx.NamedStmt
	FindByXID *sqlx.Stmt
}

func NewApplication(db *nsql.Database) *Application {
	tableName := "Application"
	columns := `"metadata","createdAt","updatedAt","modifiedBy","version","xid","name"`
	allColumns := `"id",` + columns
	namedColumns := `:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:xid,:name`

	return &Application{
		Insert:    db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
		FindByXID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "xid" = $1`, allColumns, tableName),
	}
}
