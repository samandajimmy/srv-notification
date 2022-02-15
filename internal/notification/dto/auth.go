package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
)

type Subject struct {
	Id          string
	RefId       int64
	Role        string
	FullName    string
	SubjectType constant.SubjectType
	SessionID   int64
	Metadata    map[string]string
}

func (s *Subject) ModifiedBy() *Modifier {
	return &Modifier{
		ID:       s.Id,
		Role:     s.Role,
		FullName: s.FullName,
	}
}

type ClientCredential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
