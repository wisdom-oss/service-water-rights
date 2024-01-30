package types

// DamTarget represents a target value for a dam. It contains three Quantity fields:
// - Default: the default target value for the dam
// - Steady: the target value for maintaining a steady water level in the dam
// - Max: the maximum target value for the dam
//
// The Quantity type is used to represent the value and unit of measurement for
// the target values. It consists of two fields:
//   - Value: the numeric value of the target
//   - Unit: the unit of measurement for the target value
type DamTarget struct {
	Default *Quantity `json:"default,omitempty"`
	Steady  *Quantity `json:"steady,omitempty"`
	Max     *Quantity `json:"max,omitempty"`
}
