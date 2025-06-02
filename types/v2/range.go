package v2

type Range[T any] struct {
	Lower T `json:"lower"`
	Upper T `json:"upper"`
}
