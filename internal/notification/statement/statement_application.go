package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type Application struct {
	Insert     *sqlx.NamedStmt
	FindByXID  *sqlx.Stmt
	DeleteByID *sqlx.Stmt
	Update     *sqlx.NamedStmt
}

func NewApplication(db *nsql.Database) *Application {
	tableName := `Application`
	columns := `"metadata","createdAt","updatedAt","modifiedBy","version","xid","name","apiKey"`
	allColumns := `"id",` + columns
	namedColumns := `:metadata,:createdAt,:updatedAt,:modifiedBy,:version,:xid,:name,:apiKey`
	updateNamedColumns := `metadata = :metadata, "createdAt" = :createdAt, "updatedAt" = :updatedAt, "modifiedBy" = :modifiedBy, version = :version, xid = :xid, name = :name`

	return &Application{
		Insert:     db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
		FindByXID:  db.PrepareFmt(`SELECT %s FROM "%s" WHERE "xid" = $1`, allColumns, tableName),
		DeleteByID: db.PrepareFmt(`DELETE FROM "%s" WHERE id = $1`, tableName),
		Update:     db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE id = :id`, tableName, updateNamedColumns),
	}
}
