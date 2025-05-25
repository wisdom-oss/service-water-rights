package v2

type DamTarget struct {
	Default *Quantity `json:"default"`
	Steady  *Quantity `json:"steady"`
	Max     *Quantity `json:"max"`
}
