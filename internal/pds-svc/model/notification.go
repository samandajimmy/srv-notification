package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type Notification struct {
	ID             uuid.UUID       `db:"id"`
	UserRefId      int64           `db:"userRefId"`
	Title          string          `db:"title"`
	Content        string          `db:"content"`
	ContentShort   string          `db:"contentShort"`
	ContentEncoded string          `db:"contentEncoded"`
	IsRead         bool            `db:"isRead"`
	ReadAt         sql.NullTime    `db:"readAt"`
	Metadata       json.RawMessage `db:"metadata"`
	ItemMetadata
}

type AdditionalButton struct {
	ButtonLabel   string `db:"buttonLabel"`
	TransactionId string `db:"transactionId"`
	ScreenName    string `db:"screenName"`
}

func (m *AdditionalButton) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *AdditionalButton) Value() (driver.Value, error) {
	return json.Marshal(m)
}
