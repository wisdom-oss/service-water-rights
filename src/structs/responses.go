package structs

import (
	"github.com/jackc/pgtype"
	geojson "github.com/paulmach/go.geojson"
)

type UsageLocation struct {
	// ID is the internal ID of the water right
	ID int `json:"id"`
	// WaterRight is the water right number used by Cadenza referenced with the location
	WaterRight int `json:"water_right"`
	// IsActive reflects the status of the water right location
	IsActive *bool `json:"active"`
	// IsReal reflects the reality of the water right location
	IsReal *bool `json:"real"`
	// Name contains the name for the usage location
	Name *string `json:"name"`
	// Location contains the GeoJSON representation of the usage location
	Location *geojson.Geometry `json:"geojson"`
}

type IntervalRate struct {
	Amount   *int    `json:"amount"`
	Unit     *string `json:"unit"`
	Duration *string `json:"duration"`
}

type NumericKeyedName struct {
	Key  interface{} `json:"key"`
	Name interface{} `json:"name"`
}

type Rate struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

type RawWaterRight struct {
	Id                   *int              `json:"id"`
	No                   *int              `json:"no"`
	ExternalId           *string           `json:"externalId" db:"ext_id"`
	FileReference        *string           `json:"fileReference" db:"file_ref"`
	LegalTitle           *string           `json:"legalTitle" db:"legal_title"`
	State                *string           `json:"state"`
	Subject              *string           `json:"subject"`
	Address              *string           `json:"address"`
	Annotation           *string           `json:"annotation"`
	Bailee               *string           `json:"bailee"`
	DateOfChange         *pgtype.Timestamp `json:"dateOfChange" db:"date_of_change"`
	Valid                *pgtype.Daterange `json:"valid" db:"valid"`
	GrantingAuthority    *string           `json:"grantingAuthority" db:"granting_authority"`
	RegisteringAuthority *string           `json:"registeringAuthority" db:"registering_authority"`
	WaterAuthority       *string           `json:"waterAuthority" db:"water_authority"`
}

func (r RawWaterRight) ToDetailedWaterRight() DetailedWaterRight {
	var doC *string
	if r.DateOfChange == nil {
		doC = nil
	} else {
		s := r.DateOfChange.Time.String()
		doC = &s
	}

	var validLower, validUpper *string
	if r.Valid == nil {
		validLower = nil
		validUpper = nil
	} else {
		validLowerString := r.Valid.Lower.Time.String()[:10]
		validLower = &validLowerString
		validUpperString := r.Valid.Upper.Time.String()[:10]
		validUpper = &validUpperString
	}
	return DetailedWaterRight{
		Id:            r.Id,
		No:            r.No,
		ExternalId:    r.ExternalId,
		FileReference: r.FileReference,
		LegalTitle:    r.LegalTitle,
		State:         r.State,
		Subject:       r.Subject,
		Address:       r.Address,
		Annotation:    r.Annotation,
		Bailee:        r.Bailee,
		DateOfChange:  doC,
		Valid: struct {
			Lower *string `json:"lower"`
			Upper *string `json:"upper"`
		}{
			Lower: validLower,
			Upper: validUpper,
		},
		GrantingAuthority:    r.GrantingAuthority,
		RegisteringAuthority: r.RegisteringAuthority,
		WaterAuthority:       r.WaterAuthority,
		Locations:            nil,
	}
}

type DetailedWaterRight struct {
	Id            *int    `json:"id"`
	No            *int    `json:"no"`
	ExternalId    *string `json:"externalId"`
	FileReference *string `json:"fileReference"`
	LegalTitle    *string `json:"legalTitle"`
	State         *string `json:"state"`
	Subject       *string `json:"subject"`
	Address       *string `json:"address"`
	Annotation    *string `json:"annotation"`
	Bailee        *string `json:"bailee"`
	DateOfChange  *string `json:"dateOfChange"`
	Valid         struct {
		Lower *string `json:"lower"`
		Upper *string `json:"upper"`
	} `json:"valid"`
	GrantingAuthority    *string                  `json:"grantingAuthority"`
	RegisteringAuthority *string                  `json:"registeringAuthority"`
	WaterAuthority       *string                  `json:"waterAuthority"`
	Locations            *[]DetailedUsageLocation `json:"locations"`
}

type DetailedUsageLocation struct {
	Id                     *int             `json:"id"`
	WaterRight             *int             `json:"waterRight" db:"water_right"`
	Name                   *string          `json:"name"`
	No                     *int             `json:"no"`
	Active                 *pgtype.Bool     `json:"active"`
	Location               geojson.Geometry `json:"location"`
	BasinNo                pgtype.JSON      `json:"basinNo" db:"basin_no"`
	County                 *string          `json:"county"`
	EuSurveyArea           pgtype.JSON      `json:"euSurveyArea" db:"eu_survey_area"`
	Field                  *int             `json:"field"`
	GroundwaterVolume      *string          `json:"groundwaterVolume" db:"groundwater_volume"`
	LegalScope             *string          `json:"legalScope" db:"legal_scope"`
	LocalSubDistrict       *string          `json:"localSubDistrict" db:"local_sub_district"`
	MaintenanceAssociation pgtype.JSON      `json:"maintenanceAssociation" db:"maintenance_association"`
	MunicipalArea          pgtype.JSON      `json:"municipalArea" db:"municipal_area"`
	Plot                   *string          `json:"plot"`
	Real                   *bool            `json:"real"`
	Rivershed              *string          `json:"rivershed"`
	SerialNo               *string          `json:"serialNo" db:"serial_no"`
	TopMap1To25000         pgtype.JSON      `json:"topMap1to25000" db:"top_map_1_25000"`
	WaterBody              *string          `json:"waterBody" db:"water_body"`
	FloodArea              *string          `json:"floodArea" db:"flood_area"`
	WaterProtectionArea    *string          `json:"waterProtectionArea" db:"water_protection_area"`
	WithdrawalRates        *pgtype.JSON     `json:"withdrawalRates" db:"withdrawal_rate"`
	FluidDischarge         *pgtype.JSON     `json:"fluidDischarge" db:"fluid_discharge"`
	IrrigationArea         *pgtype.JSON     `json:"irrigationArea" db:"irrigation_area"`
	RainSupplement         *pgtype.JSON     `json:"rainSupplement" db:"rain_supplement"`
}
