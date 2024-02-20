package types

import "github.com/jackc/pgx/v5/pgtype"

// NumericKeyedValue represents a key-value pair where the key is a numeric value
// and the value is a text value. This struct is typically used for storing and
// manipulating data with numeric keys.
// While the Key is always present the Value may be missing in some cases
type NumericKeyedValue struct {
	Key   *pgtype.Numeric `json:"key,omitempty"`
	Value *pgtype.Text    `json:"value,omitempty"`
}
