package v2

type InjectionLimit struct {
	Substance *string  `db:"substance" json:"substance"`
	Quantity  Quantity `db:"quantity"  json:"quantity"`
}
