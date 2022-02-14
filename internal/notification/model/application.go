package model

import (
	"database/sql"
)

type Application struct {
	ID         int64          `db:"id"`
	XID        string         `db:"xid"`
	Name       string         `db:"name"`
	ApiKey     string         `db:"apiKey"`
	WebhookURL sql.NullString `db:"webhookUrl"`
	*BaseField
}

type ApplicationListResult struct {
	Rows  []Application
	Count int64
}
