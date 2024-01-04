package types

import geojson "github.com/paulmach/go.geojson"

// UsageLocationBaseData contains all basic data every usage location has to
// have on output.
// There are still some fields with nil values to support scanning null values
// from the database to the struct
type UsageLocationBaseData struct {
	// ID contains the internal id of the usage location
	ID int `json:"id" db:"id"`

	// Name contains the name of the usage location
	Name *string `json:"name,omitempty" db:"name"`

	// WaterRightID contains the WaterRight.ID of the water right associated to
	// the usage location
	WaterRightID *int `json:"waterRight,omitempty" db:"water_right"`

	// IsActive indicates if the location is actively used
	IsActive *bool `json:"active,omitempty" db:"active"`

	// IsReal indicates if the location is real and not a virtual location
	IsReal *bool `json:"real,omitempty" db:"real"`

	// Location contains the usage location's location encoded as GeoJSON object
	Location *geojson.Geometry `json:"geojson,omitempty" db:"location"`
}
