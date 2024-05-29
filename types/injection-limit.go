package types

import "github.com/jackc/pgx/v5/pgtype"

type InjectionLimit struct {
	Substance pgtype.Text `json:"substance" db:"substance"`
	Quantity  Quantity    `json:"quantity" db:"quantity"`
}
