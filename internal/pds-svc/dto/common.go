package dto

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/constant"
)

type Modifier struct {
	ID       string                `json:"-"`
	Role     constant.ModifierRole `json:"role"`
	FullName string                `json:"full_name"`
}

type ItemMetadataResponse struct {
	CreatedAt  int64    `json:"createdAt"`
	UpdatedAt  int64    `json:"updatedAt"`
	ModifiedBy Modifier `json:"modifiedBy"`
	Version    int64    `json:"version"`
}

type FindOptions struct {
	Limit         int                    `json:"limit"`
	Skip          int                    `json:"skip"`
	SortBy        string                 `json:"sortBy"`
	SortDirection string                 `json:"sortDirection"`
	Filters       map[string]interface{} `json:"-"`
}

type ListMetadata struct {
	Count int64 `json:"count"`
	FindOptions
}
