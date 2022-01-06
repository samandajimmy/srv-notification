package model

import "encoding/json"

type ClientConfig struct {
	XID           string          `db:"xid"`
	Key           string          `db:"key"`
	Value         json.RawMessage `db:"value"`
	ApplicationId int             `db:"applicationId"`
	Metadata      json.RawMessage `db:"metadata"`
	ItemMetadata
}
