package v2

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-chrono/chrono"
	"github.com/jackc/pgx/v5/pgtype"
)

const NanosPerYear = 365 * 24 * time.Hour

type Rate struct {
	Value *float64        `db:"value" json:"amount"`
	Unit  *string         `db:"unit"  json:"unit"`
	Per   pgtype.Interval `db:"per"   json:"per"`
}

type rate struct {
	Value *float64 `json:"amount"`
	Unit  *string  `json:"unit"`
	Per   string   `json:"per"`
}

func (r Rate) MarshalJSON() ([]byte, error) {
	out := rate{
		Value: r.Value,
		Unit:  r.Unit,
	}

	period := chrono.Period{
		Years:  float32(r.Per.Months / 12),                         //nolint:mnd
		Months: float32(r.Per.Months - ((r.Per.Months / 12) * 12)), //nolint:mnd
		Days:   float32(r.Per.Days),
	}

	duration := chrono.DurationOf(chrono.Extent(r.Per.Microseconds * 1000)) //nolint:mnd

	out.Per = chrono.FormatDuration(period, duration)
	fmt.Println(out.Per)
	return json.Marshal(out)
}
