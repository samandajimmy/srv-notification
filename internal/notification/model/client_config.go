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
	*BaseField
}

type ClientConfigDetailed struct {
	ClientConfig *ClientConfig `db:"cc"`
	Application  *Application  `db:"a"`
}

type ClientConfigListResult struct {
	Rows  []ClientConfigDetailed
	Count int64
}

type SMTPConfig struct {
	Host     string `json:"SMTP_HOST"`
	Port     string `json:"SMTP_PORT"`
	Username string `json:"SMTP_USERNAME"`
	Password string `json:"SMTP_PASSWORD"`
}
