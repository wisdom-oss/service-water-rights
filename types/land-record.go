package types

import "github.com/jackc/pgx/v5/pgtype"

type LandRecord struct {
	Fallback *pgtype.Text `json:"fallback,omitempty"`
	District *pgtype.Text `json:"district,omitempty"`
	Field    *pgtype.Int8 `json:"field,omitempty"`
}
