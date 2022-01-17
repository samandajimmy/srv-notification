package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type ClientConfig struct {
	FindByKey *sqlx.Stmt
	Insert    *sqlx.NamedStmt
}

func NewClientConfig(db *nsql.Database) *ClientConfig {
	tableName := `ClientConfig`
	columns := `"xid","key","value","applicationId","metadata","createdAt","updatedAt","modifiedBy","version"`
	namedColumns := `:createdAt,:updatedAt,:metadata,:modifiedBy,:version,:key,:value,:applicationId,:xid`

	return &ClientConfig{
		FindByKey: db.PrepareFmt(`SELECT %s FROM "%s" WHERE "key" = $1 AND "applicationId" = $2`, columns, tableName),
		Insert:    db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
	}
}
