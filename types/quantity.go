package types

// Quantity represents a numerical quantity with a unit
type Quantity struct {
	Value *float64 `json:"value,omitempty"`
	Unit  *string  `json:"unit,omitempty"`
}
