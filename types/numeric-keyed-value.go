package types

import "github.com/jackc/pgx/v5/pgtype"

// While the Key is always present the Value may be missing in some cases.
type NumericKeyedValue struct {
	Key   *pgtype.Numeric `json:"key,omitempty"`
	Value *pgtype.Text    `json:"value,omitempty"`
}
