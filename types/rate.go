package types

import "github.com/jackc/pgx/v5/pgtype"

// Rate represents a rate quantity and its corresponding per interval.
// The rate quantity is described by the `Quantity` struct, which consists of a
// numeric value and a text unit.
// The per interval is represented by the `Per` field of type `pgtype.Interval`.
// Together, they define the rate of change or occurrence for a given quantity.
// Refer to the `types.Quantity` struct for more information on the quantity
// representation.
// Example usage:
//
//	rate := Rate{
//	  Quantity: Quantity{
//	    Value: pgtype.Numeric{Value: 2},
//	    Unit:  pgtype.Text{String: "meters"},
//	  },
//	  Per: pgtype.Interval{Duration: time.Hour},
//	}
//	// rate represents a rate of 2 meters per hour.
type Rate struct {
	Value *pgtype.Numeric  `json:"value,omitempty" db:"value"`
	Unit  *pgtype.Text     `json:"unit,omitempty" db:"unit"`
	Per   *pgtype.Interval `json:"per,omitempty"`
}
