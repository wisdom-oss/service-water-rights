package types

import "github.com/jackc/pgx/v5/pgtype"

type InjectionLimit struct {
	Substance pgtype.Text `db:"substance" json:"substance"`
	Quantity  Quantity    `db:"quantity"  json:"quantity"`
}
