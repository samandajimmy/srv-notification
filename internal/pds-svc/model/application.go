package model

import "encoding/json"

type Application struct {
	ID       string          `db:"id"`
	XID      string          `db:"xid"`
	Name     string          `db:"name"`
	Metadata json.RawMessage `db:"metadata"`
	ItemMetadata
}
