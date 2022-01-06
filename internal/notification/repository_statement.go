package notification

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/statement"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

type RepositoryStatement struct {
	ClientConfig *statement.ClientConfigStatement
}

// NewRepositoryStatement prepare all sql statements
func NewRepositoryStatement(db *nsql.Database) *RepositoryStatement {
	rs := RepositoryStatement{
		ClientConfig: statement.NewClientConfigStatement(db),
	}
	return &rs
}
