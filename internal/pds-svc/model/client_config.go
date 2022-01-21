package model

import (
	"encoding/json"
)

type ClientConfig struct {
	ID            int64           `db:"id"`
	XID           string          `db:"xid"`
	Key           string          `db:"key"`
	Value         json.RawMessage `db:"value"`
	ApplicationId int64           `db:"applicationId"`
	Metadata      json.RawMessage `db:"metadata"`
	ItemMetadata
}

type ClientConfigVO struct {
	ID             int64           `db:"id"`
	XID            string          `db:"xid"`
	Key            string          `db:"key"`
	Value          json.RawMessage `db:"value"`
	ApplicationXid string          `db:"applicationXid"`
	Metadata       json.RawMessage `db:"metadata"`
	ItemMetadata
}

type ClientConfigSearchResult struct {
	Rows  []ClientConfigVO
	Count int64
}
