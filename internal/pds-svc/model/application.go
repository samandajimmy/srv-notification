package model

import (
	"database/sql"
	"encoding/json"
)

type Application struct {
	ID         int64           `db:"id"`
	XID        string          `db:"xid"`
	Name       string          `db:"name"`
	ApiKey     string          `db:"apiKey"`
	WebhookURL sql.NullString  `db:"webhookUrl"`
	Metadata   json.RawMessage `db:"metadata"`
	ItemMetadata
}

type ApplicationFindResult struct {
	Rows  []Application
	Count int64
}
