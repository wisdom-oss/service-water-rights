package v2

type Quantity struct {
	Value *float64 `db:"value" json:"amount"`
	Unit  *string  `db:"unit"  json:"unit"`
}
