package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type ClientConfigStatement struct {
	FindByKey *sqlx.Stmt
}

func NewClientConfigStatement(db *nsql.Database) *ClientConfigStatement {
	tableName := "ClientConfig"
	columns := `"xid","key","value","applicationId","metadata","createdAt","updatedAt","modifiedBy","version"`

	return &ClientConfigStatement{
		FindByKey: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "key" = $1 AND "applicationId" = $2`, columns, tableName),
	}
}
