package statement

import (
	"github.com/jmoiron/sqlx"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type Notification struct {
	Insert   *sqlx.NamedStmt
	FindByID *sqlx.Stmt
	Update   *sqlx.NamedStmt
}

func NewNotification(db *nsql.Database) *Notification {
	tableName := `Notification`
	columns := `"id","createdAt","updatedAt","modifiedBy","metadata","version","applicationId","userRefId","isRead","readAt","options"`
	namedColumns := `:id,:createdAt,:updatedAt,:modifiedBy,:metadata,:version,:applicationId,:userRefId,:isRead,:readAt,:options`
	updateNamedColumns := `"id" = :id,"createdAt" = :createdAt,"updatedAt" = :updatedAt,"modifiedBy" = :modifiedBy,"metadata" = :metadata,"version" = :version,"applicationId" = :applicationId,"userRefId" = :userRefId,"isRead" = :isRead,"readAt" = :readAt,"options" = :options`

	return &Notification{
		Insert:   db.PrepareNamedFmt(`INSERT INTO "%s"(%s) VALUES (%s)`, tableName, columns, namedColumns),
		FindByID: db.PrepareFmt(`SELECT %s FROM "%s" WHERE id = $1`, columns, tableName),
		Update:   db.PrepareNamedFmt(`UPDATE "%s" SET %s WHERE "id" = :id`, tableName, updateNamedColumns),
	}
}
