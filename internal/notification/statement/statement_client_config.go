package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type ClientConfig struct {
	FindByKey  *sqlx.Stmt
	FindByXID  *sqlx.Stmt
	Insert     *sqlx.NamedStmt
	UpdateByID *sqlx.NamedStmt
	DeleteByID *sqlx.Stmt
}

func NewClientConfig(db *nsql.Database) *ClientConfig {
	tableName := `ClientConfig`
	columns := `"createdAt","updatedAt","metadata","modifiedBy","version","key","value","applicationId","xid"`
	namedColumns := `:createdAt,:updatedAt,:metadata,:modifiedBy,:version,:key,:value,:applicationId,:xid`
	allColumns := `"id",` + columns
	updatedNamedColumns := `"updatedAt" = :updatedAt, "metadata" = :metadata, "modifiedBy" = :modifiedBy, "version" = :version, "key" = :key, "value" = :value, "applicationId" = :applicationId, "xid" = :xid`

	return &ClientConfig{
		FindByKey:  db.PrepareFmt(`SELECT %s FROM "%s" WHERE "key" = $1 AND "applicationId" = $2`, columns, tableName),
		Insert:     db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
		FindByXID:  db.PrepareFmt(`SELECT %s FROM "%s" WHERE "xid" = $1`, allColumns, tableName),
		UpdateByID: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "id" = :id`, tableName, updatedNamedColumns),
		DeleteByID: db.PrepareFmt(`DELETE FROM "%s" WHERE "id" = $1`, tableName),
	}
}
