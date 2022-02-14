package dto

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
)

type Modifier struct {
	ID       string                `json:"-"`
	Role     constant.ModifierRole `json:"role"`
	FullName string                `json:"fullName"`
}

type BaseField struct {
	CreatedAt  int64     `json:"createdAt"`
	UpdatedAt  int64     `json:"updatedAt"`
	ModifiedBy *Modifier `json:"modifiedBy"`
	Version    int64     `json:"version"`
}

type ListPayload struct {
	Limit   int64             `json:"limit" query:"limit"`
	Skip    int64             `json:"skip" query:"skip"`
	SortBy  string            `json:"sortBy" query:"sortBy"`
	Filters map[string]string `json:"-" query:"filters"`
}

type ListMetadata struct {
	Count  int64  `json:"count"`
	Limit  int64  `json:"limit"`
	Skip   int64  `json:"skip"`
	SortBy string `json:"sortBy"`
}

func ToListMetadata(p *ListPayload, count int64) *ListMetadata {
	return &ListMetadata{
		Count:  count,
		Limit:  p.Limit,
		Skip:   p.Skip,
		SortBy: p.SortBy,
	}
}
