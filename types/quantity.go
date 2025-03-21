package types

import "github.com/jackc/pgx/v5/pgtype"

// Quantity represents a numerical quantity with a unit.
type Quantity struct {
	Value pgtype.Numeric `db:"value" json:"value,omitempty"`
	Unit  pgtype.Text    `db:"unit"  json:"unit,omitempty"`
}
