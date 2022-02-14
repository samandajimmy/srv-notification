package model

import (
	"database/sql/driver"
	"encoding/json"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"time"
)

type Modifier struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"fullName"`
}

func (m *Modifier) Scan(src interface{}) error {
	return nsql.ScanJSON(src, m)
}

func (m *Modifier) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func ToModifierDTO(model *Modifier) *dto.Modifier {
	return &dto.Modifier{
		ID:       model.ID,
		Role:     model.Role,
		FullName: model.FullName,
	}
}

func ToModifier(dto *dto.Modifier) *Modifier {
	return &Modifier{
		ID:       dto.ID,
		Role:     dto.Role,
		FullName: dto.FullName,
	}
}

type BaseField struct {
	CreatedAt  time.Time       `db:"createdAt"`
	UpdatedAt  time.Time       `db:"updatedAt"`
	ModifiedBy *Modifier       `db:"modifiedBy"`
	Version    int64           `db:"version"`
	Metadata   json.RawMessage `db:"metadata"`
}

func (m *BaseField) Upgrade(modifiedBy Modifier, opts ...time.Time) *BaseField {
	var t time.Time
	if len(opts) > 0 {
		t = opts[0]
	} else {
		t = time.Now()
	}

	return &BaseField{
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  t,
		ModifiedBy: &modifiedBy,
		Version:    m.Version + 1,
	}
}

func NewBaseField(modifiedBy *Modifier) *BaseField {
	// Init timestamp
	t := time.Now()

	return &BaseField{
		CreatedAt:  t,
		UpdatedAt:  t,
		ModifiedBy: modifiedBy,
		Version:    1,
		Metadata:   []byte("{}"),
	}
}

func ToBaseFieldDTO(m *BaseField) *dto.BaseField {
	return &dto.BaseField{
		UpdatedAt:  m.UpdatedAt.Unix(),
		CreatedAt:  m.CreatedAt.Unix(),
		ModifiedBy: ToModifierDTO(m.ModifiedBy),
		Version:    m.Version,
	}
}
