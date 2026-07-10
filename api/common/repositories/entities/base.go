package entities

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type BaseInt struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type BaseNoID struct {
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type QueryParams struct {
	Sort   string        `json:"sort"`
	Limit  int           `json:"limit,string"`
	Page   int           `json:"page,string"`
	Search []SearchField `json:"search"`
	Query  Query         `json:"query"`

	SelectFields []string `json:"selectFields,omitempty"`

	// Fields to update or increment
	UpdateFields    []string `json:"updateFields"`
	IncrementFields []string `json:"incrementFields,omitempty"`

	// Fields to check for conflicts when inserting or updating
	ConflictColumns []string `json:"conflictColumns,omitempty"`
}

type Query struct {
	Filters    string `json:"filters"`
	Joins      string `json:"joins"`
	JoinValues []any  `json:"joinValues,omitempty"`
	Values     []any  `json:"values"`
	From       string `json:"from,omitempty"` // << NEW: tabelas/aliases para UPDATE/DELETE no Postgres
}

type SearchField struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type PaginatedResult struct {
	Data         interface{} `json:"data"`
	TotalItems   int64       `json:"totalItems"`
	TotalPages   int64       `json:"totalPages"`
	CurrentPage  int64       `json:"currentPage"`
	PreviousPage int64       `json:"previousPage"`
	NextPage     int64       `json:"nextPage"`
}

type JSONB map[string]any

func (j *JSONB) Cast(data any) error {
	switch v := data.(type) {
	case JSONB:
		*j = v

	case map[string]any:
		*j = JSONB(v)

	case []byte:
		return json.Unmarshal(v, j)

	case string:
		return json.Unmarshal([]byte(v), j)

	default:
		return fmt.Errorf("invalid reward data type: %T", data)
	}

	return nil
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(src any) error {
	if src == nil {
		*j = nil
		return nil
	}

	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("JSONB.Scan: unsupported source type: %T", src)
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber() // evita float64; vira json.Number

	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return err
	}

	*j = m
	return nil
}
