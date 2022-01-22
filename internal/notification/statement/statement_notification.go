package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type Notification struct {
	Insert *sqlx.NamedStmt
}

func NewNotification(db *nsql.Database) *Notification {
	tableName := `Notification`
	columns := `"id","createdAt","updatedAt","modifiedBy","metadata","version","userRefId","title","content","contentShort","contentEncoded","isRead","readAt"`
	namedColumns := `:id,:createdAt,:updatedAt,:modifiedBy,:metadata,:version,:userRefId,:title,:content,:contentShort,:contentEncoded,:isRead,:readAt`

	return &Notification{
		Insert: db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
	}
}
