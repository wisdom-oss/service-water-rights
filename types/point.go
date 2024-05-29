package types

import (
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type Location struct{ ewkb.Point }

func (l Location) MarshalJSON() ([]byte, error) {
	return geojson.Marshal(l.Point.Point)
}
