package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type ClientConfig struct {
	FindByKey *sqlx.Stmt
}

func NewClientConfig(db *nsql.Database) *ClientConfig {
	tableName := "ClientConfig"
	columns := `"xid","key","value","applicationId","metadata","createdAt","updatedAt","modifiedBy","version"`

	return &ClientConfig{
		FindByKey: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "key" = $1 AND "applicationId" = $2`, columns, tableName),
	}
}
