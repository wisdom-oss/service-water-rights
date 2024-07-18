package types

import (
	"encoding/json"

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
	Value *pgtype.Numeric  `json:"value,omitempty" db:"value"`
	Unit  *pgtype.Text     `json:"unit,omitempty" db:"unit"`
	Per   *pgtype.Interval `json:"per,omitempty"`
}

const (
	microsecondsPerSecond int64 = 1000000
	microsecondsPerMinute int64 = 60 * microsecondsPerSecond
	microsecondsPerHour         = 60 * microsecondsPerMinute
	microsecondsPerDay          = 24 * microsecondsPerHour
	microsecondsPerMonth        = 30 * microsecondsPerDay
	microsecondsPerYear         = 365 * microsecondsPerDay
)

const litrePerCubicMeter = 1000

func (r Rate) CubicMeterPerYear() float64 {
	f64, err := r.Value.Float64Value()
	if err != nil {
		return 0
	}
	amount := f64.Float64
	totalMicros := r.Per.Microseconds + int64(r.Per.Days)*microsecondsPerDay + int64(r.Per.Months)*microsecondsPerMonth

	microRelation := float64(totalMicros) / float64(microsecondsPerYear)

	yearlyAmount := amount / microRelation

	switch r.Unit.String {
	case "l":
		return yearlyAmount / litrePerCubicMeter
	case "m³":
		return yearlyAmount
	default:
		return 0
	}
}

func (r Rate) MarshalJSON() ([]byte, error) {
	type outputRate struct {
		Value *pgtype.Numeric `json:"value,omitempty" db:"value"`
		Unit  *pgtype.Text    `json:"unit,omitempty" db:"unit"`
		Per   int64           `json:"per,omitempty"`
	}
	totalMicros := r.Per.Microseconds + int64(r.Per.Days)*microsecondsPerDay + int64(r.Per.Months)*microsecondsPerMonth
	or := outputRate{
		Value: r.Value,
		Unit:  r.Unit,
		Per:   totalMicros,
	}
	return json.Marshal(or)
}
