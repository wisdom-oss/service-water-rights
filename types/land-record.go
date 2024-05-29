package types

import "github.com/jackc/pgx/v5/pgtype"

type LandRecord struct {
	District *pgtype.Text `json:"district,omitempty"`
	Field    *pgtype.Int8 `json:"field,omitempty"`
	Fallback *pgtype.Text `json:"fallback,omitempty"`
}
