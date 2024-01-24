package types

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UsageLocation struct {
	// ID identifies the usage location
	ID pgtype.Int8 `json:"id,omitempty" db:"id"`

	// SerialID enumerates the usage location inside a water right
	SerialID *pgtype.Text `json:"serial,omitempty" db:"no"`

	// WaterRightID shows which water right is associated with this usage
	// location
	WaterRightID *pgtype.Int8 `json:"waterRight,omitempty" db:"water_right"`

	// Active shows if the usage location is currently used
	Active *pgtype.Bool `json:"active,omitempty" db:"active"`

	// Real shows if the usage location actually exists or not
	Real *pgtype.Bool `json:"real,omitempty" db:"real"`

	// Name either contains the usage locations name or a different descriptor
	Name *pgtype.Text `json:"name,omitempty" db:"name"`

	// LegalPurpose contains the legal purpose for the usage location
	LegalPurpose *[2]string `json:"legalPurpose,omitempty" db:"legal_purpose"`

	// MapExcerpt contains the identification of the area the location is in on
	// a topological map using a 1:25000 scale
	MapExcerpt *NumericKeyedValue `json:"mapExcerpt,omitempty" db:"map_excerpt"`

	// MunicipalArea contains the ARS of the municipal and the name in which
	// this usage location is located
	MunicipalArea *NumericKeyedValue `json:"municipalArea,omitempty" db:"municipal_area"`

	// County contains the name of the county the usage location is located in
	County *pgtype.Text `json:"county,omitempty" db:"county"`

	// LandRecord contains information about the parish the usage location is
	// located in
	LandRecord *LandRecord `json:"landRecord,omitempty" db:"land_record"`

	// Plot contains information about the plot the usage location is located
	// in
	Plot *pgtype.Text `json:"plot,omitempty" db:"plot"`

	// MaintenanceAssociation contains information about the (legal) person
	// responsible for maintaining the usage location
	MaintenanceAssociation *NumericKeyedValue `json:"maintenanceAssociation,omitempty" db:"maintenance_association"`

	// EUSurveyArea contains information about an EU-specified survey area
	EUSurveyArea *NumericKeyedValue `json:"euSurveyArea,omitempty" db:"eu_survey_area"`

	// CatchmentAreaCode contains further information about the area the usage
	// location is located in
	CatchmentAreaCode *NumericKeyedValue `json:"catchmentAreaCode,omitempty" db:"catchment_area_code"`

	// RegulationCitation contains a citation from the regulations about water
	// rights
	RegulationCitation *pgtype.Text `json:"regulationCitation,omitempty" db:"regulation_citation"`

	// WithdrawalRates contains information about the rates of water withdrawal
	// that are reflected by the associated water right
	WithdrawalRates []Rate `json:"withdrawalRates,omitempty" db:"withdrawal_rates"`

	// PumpingRates contains information about the allowed pumping rates for the
	// usage location
	PumpingRates []Rate `json:"pumpingRates,omitempty" db:"pumping_rates"`

	// InjectionRates contains information about the allowed rates for
	// injecting something at the usage location
	InjectionRates []Rate `json:"injectionRates,omitempty" db:"injection_rates"`

	// WasteWaterFlowVolume contains information about the amount of waste water
	// handled at that location
	WasteWaterFlowVolume []Rate `json:"wasteWaterFlowVolume,omitempty" db:"waste_water_flow_volume"`

	// RiverBasin describes the river used in connection with the location
	RiverBasin *pgtype.Text `json:"riverBasin,omitempty" db:"river_basin"`

	// GroundwaterBody describes the groundwater used at this location
	GroundwaterBody *pgtype.Text `json:"groundwaterBody,omitempty" db:"groundwater_body"`

	// WaterBody describes the water body used at this location
	WaterBody *pgtype.Text `json:"waterBody,omitempty" db:"water_body"`

	// FloodArea describes the area in which a flood may happen
	FloodArea *pgtype.Text `json:"floodArea,omitempty" db:"flood_area"`

	// WaterProtectionArea contains information about a possible water
	// protection area
	WaterProtectionArea *pgtype.Text `json:"waterProtectionArea,omitempty" db:"water_protection_area"`

	// DamTargetLevels contains information about the levels that are to be
	// expected at a dam
	DamTargetLevels *DamTarget `json:"damTargetLevels,omitempty" db:"dam_target_levels"`

	// FluidDischarge contains rates of discharged fluids
	FluidDischarge []Rate `json:"fluidDischarge,omitempty" db:"fluid_discharge"`

	// RainSupplement contains information about additionally used rain water
	RainSupplement []Rate `json:"rainSupplement,omitempty" db:"rain_supplement"`

	// IrrigationArea contains information about the area that the rain is
	// collected from
	IrrigationArea *Quantity `json:"irrigationArea,omitempty" db:"irrigation_area"`

	// PHValues contains information about the ph values at the usage location
	PHValues *pgtype.Range[pgtype.Numeric] `json:"phValues,omitempty" db:"ph_values"`

	// InjectionLimits contains information about injection limitations
	InjectionLimits [][2]interface{} `json:"injectionLimits,omitempty" db:"injection_limit"`

	// Location contains the GeoJSON representation of the usage locations
	// location
	Location interface{} `json:"location,omitempty" db:"location"`
}
