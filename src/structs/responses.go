package structs

import (
	geojson "github.com/paulmach/go.geojson"
)

// UsageLocation contains some basic information as the location and the associated
// water right and the status of the water right
type UsageLocation struct {
	// ID is the internal ID of the water right
	ID int `json:"id"`
	// Name contains the name for the usage location
	Name *string `json:"name"`
	// WaterRight is the water right number used by Cadenza referenced with the location
	WaterRight int `json:"waterRight"`
	// IsActive reflects the status of the water right location
	IsActive *bool `json:"active"`
	// IsReal reflects the reality of the water right location
	IsReal *bool `json:"real"`
	// Location contains the GeoJSON representation of the usage location
	Location *geojson.Geometry `json:"geojson"`
}

// DetailedUsageLocation contains more information about a usage location that
// may be stored in the database. Due to the parsing process of the water rights
// all fields are pointers to allow the usage of nil on every field to mark the
// absence of the data in the database
type DetailedUsageLocation struct {
	// use the already defined UsageLocation as a base and extend it with the
	// more detailed fields
	UsageLocation
	NlwknID                *int              `json:"no"`
	BasinNumber            *NumericKeyedName `json:"basinNo"`
	County                 *string           `json:"county"`
	EuSurveyArea           *NumericKeyedName `json:"euSurveyArea"`
	Field                  *int              `json:"field"`
	GroundwaterVolume      *string           `json:"groundwaterVolume"`
	LegalScope             *string           `json:"legalScope"`
	LocalSubDistrict       *string           `json:"localSubDistrict"`
	MaintenanceAssociation *NumericKeyedName `json:"maintenanceAssociation"`
	MunicipalArea          *NumericKeyedName `json:"municipalArea"`
	Plot                   *string           `json:"plot"`
	Rivershed              *string           `json:"rivershed"`
	SerialNumber           *string           `json:"serialNo"`
	TopMap1To25000         *NumericKeyedName `json:"topMap1To25000"`
	WaterBody              *string           `json:"waterBody"`
	FloodArea              *string           `json:"floodArea"`
	WaterProtectionArea    *string           `json:"waterProtectionArea"`
	WithdrawalRates        []*IntervalRate   `json:"withdrawalRates"`
	FluidDischarge         []*IntervalRate   `json:"fluidDischarge"`
	IrrigationArea         *Rate             `json:"irrigationArea"`
	RainSupplement         []*IntervalRate   `json:"rainSupplement"`
}

type WaterRight struct {
	ID                   int            `json:"id"`
	NlwknID              int            `json:"no"`
	ExternalId           *string        `json:"externalId"`
	FileReference        *string        `json:"fileReference"`
	LegalTitle           *string        `json:"legalTitle"`
	State                *string        `json:"state"`
	Subject              *string        `json:"subject"`
	Address              *string        `json:"address"`
	Annotation           *string        `json:"annotation"`
	Bailee               *string        `json:"bailee"`
	DateOfChange         *int64         `json:"dateOfChange"`
	Valid                *DateTimeRange `json:"valid"`
	GrantingAuthority    *string        `json:"grantingAuthority"`
	RegisteringAuthority *string        `json:"registeringAuthority"`
	WaterAuthority       *string        `json:"waterAuthority"`
}

type WaterRightDetailResponse struct {
	WaterRight
	Locations []DetailedUsageLocation `json:"locations"`
}
