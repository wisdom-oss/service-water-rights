package structs

import geojson "github.com/paulmach/go.geojson"

type UsageLocation struct {
	// ID is the internal ID of the water right
	ID int `json:"id"`
	// WaterRight is the water right number used by Cadenza referenced with the location
	WaterRight int `json:"water_right"`
	// IsActive reflects the status of the water right location
	IsActive bool `json:"active"`
	// IsReal reflects the reality of the water right location
	IsReal bool `json:"real"`
	// Name contains the name for the usage location
	Name string `json:"name"`
	// Location contains the GeoJSON representation of the usage location
	Location geojson.Geometry `json:"geojson"`
}
