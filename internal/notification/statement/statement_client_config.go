package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type ClientConfig struct {
	FindByKey                *sqlx.Stmt
	FindByXID                *sqlx.Stmt
	Insert                   *sqlx.NamedStmt
	UpdateByID               *sqlx.NamedStmt
	DeleteByID               *sqlx.Stmt
	FindJoinApplicationByXID *sqlx.Stmt
}

func NewClientConfig(db *nsql.Database) *ClientConfig {
	tableName := `ClientConfig`
	columns := `"createdAt","updatedAt","metadata","modifiedBy","version","key","value","applicationId","xid"`
	namedColumns := `:createdAt,:updatedAt,:metadata,:modifiedBy,:version,:key,:value,:applicationId,:xid`
	allColumns := `"id",` + columns
	updatedNamedColumns := `"updatedAt" = :updatedAt, "metadata" = :metadata, "modifiedBy" = :modifiedBy, "version" = :version, "key" = :key, "value" = :value, "applicationId" = :applicationId, "xid" = :xid`

	applicationColumn := `"Application"."xid" AS "applicationXid"`
	columnsWithExplicitTable := `"ClientConfig"."createdAt", "ClientConfig"."updatedAt", "ClientConfig"."metadata", "ClientConfig"."modifiedBy", "ClientConfig"."version", "ClientConfig"."key", "ClientConfig"."value", "ClientConfig"."xid"`
	columnsWithExplicitTable += ", " + applicationColumn
	joinApplication := `LEFT JOIN "Application" ON "Application"."id" = "ClientConfig"."applicationId"`

	return &ClientConfig{
		FindByKey:  db.PrepareFmt(`SELECT %s FROM "%s" WHERE "key" = $1 AND "applicationId" = $2`, columns, tableName),
		Insert:     db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
		FindByXID:  db.PrepareFmt(`SELECT %s FROM "%s" WHERE "xid" = $1`, allColumns, tableName),
		UpdateByID: db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "id" = :id`, tableName, updatedNamedColumns),
		DeleteByID: db.PrepareFmt(`DELETE FROM "%s" WHERE "id" = $1`, tableName),
		FindJoinApplicationByXID: db.PrepareFmt(`SELECT %s FROM "%s" %s WHERE "%s"."xid" = $1`,
			columnsWithExplicitTable, tableName, joinApplication, tableName,
		),
	}
}
