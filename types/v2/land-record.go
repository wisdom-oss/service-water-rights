package v2

type LandRecord struct {
	District *string `json:"district"`
	Field    *int64  `json:"field"`
	Fallback *string `json:"fallback"`
}
