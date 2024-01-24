package types

import "github.com/jackc/pgx/v5/pgtype"

type LandRecord struct {
	RegisteringDistrict *pgtype.Text `json:"registeringDistrict,omitempty"`
	FieldNumber         *pgtype.Int8 `json:"fieldNumber,omitempty"`
}
