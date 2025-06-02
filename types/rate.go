package types

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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
	Value *pgtype.Numeric  `db:"value"           json:"value,omitempty"`
	Unit  *pgtype.Text     `db:"unit"            json:"unit,omitempty"`
	Per   *pgtype.Interval `json:"per,omitempty"`
}

const (
	oneDay   = 24 * time.Hour
	oneMonth = 30 * oneDay
	oneYear  = 12 * oneMonth
)

func (r Rate) CubicMeterPerYear() float64 {
	fl64, err := r.Value.Float64Value()
	if err != nil {
		return 0
	}

	amount := fl64.Float64
	micros := r.Per.Microseconds + int64(r.Per.Days)*oneDay.Microseconds() + int64(r.Per.Months)

	relativeMicros := float64(micros) / float64(oneYear.Microseconds())

	yearlyAmount := amount / relativeMicros

	switch r.Unit.String {
	case "l", "L", "liter", "litre", "Liter", "Litre":
		return yearlyAmount / 1000 //nolint:mnd
	case "m³", "m^3", "m3":
		return yearlyAmount
	default:
		return 0
	}
}
