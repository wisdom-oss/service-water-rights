package types

import "github.com/jackc/pgx/v5/pgtype"

// Quantity represents a numerical quantity with a unit
type Quantity struct {
	Value pgtype.Numeric `json:"value,omitempty" db:"value"`
	Unit  pgtype.Text    `json:"unit,omitempty" db:"unit"`
}
