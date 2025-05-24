package types

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/wroge/wgs84/v2"
)

type UsageLocation struct {
	// ID identifies the usage location in the database
	ID pgtype.Int8 `db:"id" json:"id,omitempty"`

	// LocationNumber is the usage locations id issued by the cadenza database
	LocationNumber *pgtype.Int8 `db:"no" json:"no"`

	// SerialID enumerates the usage location inside water right
	SerialID *pgtype.Text `db:"serial" json:"serial,omitempty"`

	// WaterRightID shows which water right is associated with this usage
	// location
	WaterRightID pgtype.Int8 `db:"water_right" json:"waterRight,omitempty"`

	// LegalDepartment shows into which legal department this usage location falls
	LegalDepartment pgtype.Text `db:"legal_department" json:"legalDepartment"`

	// Active shows if the usage location is currently used
	Active *pgtype.Bool `db:"active" json:"active,omitempty"`

	// Real shows if the usage location actually exists or not
	Real *pgtype.Bool `db:"real" json:"real,omitempty"`

	// Name either contains the usage locations name or a different descriptor
	Name *pgtype.Text `db:"name" json:"name,omitempty"`

	// LegalPurpose contains the legal purpose for the usage location
	LegalPurpose *[]string `db:"legal_purpose" json:"legalPurpose,omitempty"`

	// MapExcerpt contains the identification of the area the location is in on
	// a topological map using a 1:25000 scale
	MapExcerpt *NumericKeyedValue `db:"map_excerpt" json:"mapExcerpt,omitempty"`

	// MunicipalArea contains the ARS of the municipal and the name in which
	// this usage location is located
	MunicipalArea *NumericKeyedValue `db:"municipal_area" json:"municipalArea,omitempty"`

	// County contains the name of the county the usage location is located in
	County *pgtype.Text `db:"county" json:"county,omitempty"`

	// LandRecord contains information about the parish the usage location is
	// located in
	LandRecord *LandRecord `db:"land_record" json:"landRecord,omitempty"`

	// Plot contains information about the plot the usage location is located
	// in
	Plot *pgtype.Text `db:"plot" json:"plot,omitempty"`

	// MaintenanceAssociation contains information about the (legal) person
	// responsible for maintaining the usage location
	MaintenanceAssociation *NumericKeyedValue `db:"maintenance_association" json:"maintenanceAssociation,omitempty"`

	// EUSurveyArea contains information about an EU-specified survey area
	EUSurveyArea *NumericKeyedValue `db:"eu_survey_area" json:"euSurveyArea,omitempty"`

	// CatchmentAreaCode contains further information about the area the usage
	// location is located in
	CatchmentAreaCode *NumericKeyedValue `db:"catchment_area_code" json:"catchmentAreaCode,omitempty"`

	// RegulationCitation contains a citation from the regulations about water
	// rights
	RegulationCitation *pgtype.Text `db:"regulation_citation" json:"regulationCitation,omitempty"`

	// WithdrawalRates contains information about the rates of water withdrawal
	// that are reflected by the associated water right
	WithdrawalRates []Rate `db:"withdrawal_rates" json:"withdrawalRates,omitempty"`

	// PumpingRates contains information about the allowed pumping rates for the
	// usage location
	PumpingRates []Rate `db:"pumping_rates" json:"pumpingRates,omitempty"`

	// InjectionRates contains information about the allowed rates for
	// injecting something at the usage location
	InjectionRates []Rate `db:"injection_rates" json:"injectionRates,omitempty"`

	// WasteWaterFlowVolume contains information about the amount of waste water
	// handled at that location
	WasteWaterFlowVolume []Rate `db:"waste_water_flow_volume" json:"wasteWaterFlowVolume,omitempty"`

	// RiverBasin describes the river used in connection with the location
	RiverBasin *pgtype.Text `db:"river_basin" json:"riverBasin,omitempty"`

	// GroundwaterBody describes the groundwater used at this location
	GroundwaterBody *pgtype.Text `db:"groundwater_body" json:"groundwaterBody,omitempty"`

	// WaterBody describes the water body used at this location
	WaterBody *pgtype.Text `db:"water_body" json:"waterBody,omitempty"`

	// FloodArea describes the area in which a flood may happen
	FloodArea *pgtype.Text `db:"flood_area" json:"floodArea,omitempty"`

	// WaterProtectionArea contains information about a possible water
	// protection area
	WaterProtectionArea *pgtype.Text `db:"water_protection_area" json:"waterProtectionArea,omitempty"`

	// DamTargetLevels contains information about the levels that are to be
	// expected at a dam
	DamTargetLevels *DamTarget `db:"dam_target_levels" json:"damTargetLevels,omitempty"`

	// FluidDischarge contains rates of discharged fluids
	FluidDischarge []Rate `db:"fluid_discharge" json:"fluidDischarge,omitempty"`

	// RainSupplement contains information about additionally used rain water
	RainSupplement []Rate `db:"rain_supplement" json:"rainSupplement,omitempty"`

	// IrrigationArea contains information about the area that the rain is
	// collected from
	IrrigationArea *Quantity `db:"irrigation_area" json:"irrigationArea,omitempty"`

	// PHValues contains information about the ph values at the usage location
	PHValues *pgtype.Range[pgtype.Numeric] `db:"ph_values" json:"phValues,omitempty"`

	// InjectionLimits contains information about injection limitations
	InjectionLimits []InjectionLimit `db:"injection_limits" json:"injectionLimits,omitempty"`

	// Location contains the GeoJSON representation of the usage locations
	// location
	Location geom.T `db:"location" json:"location"`
}

type usageLocation struct {
	// ID identifies the usage location in the database
	ID pgtype.Int8 `db:"id" json:"id,omitempty"`

	// LocationNumber is the usage locations id issued by the cadenza database
	LocationNumber *pgtype.Int8 `db:"no" json:"no"`

	// SerialID enumerates the usage location inside water right
	SerialID *pgtype.Text `db:"serial" json:"serial,omitempty"`

	// WaterRightID shows which water right is associated with this usage
	// location
	WaterRightID pgtype.Int8 `db:"water_right" json:"waterRight,omitempty"`

	// LegalDepartment shows into which legal department this usage location falls
	LegalDepartment pgtype.Text `db:"legal_department" json:"legalDepartment"`

	// Active shows if the usage location is currently used
	Active *pgtype.Bool `db:"active" json:"active,omitempty"`

	// Real shows if the usage location actually exists or not
	Real *pgtype.Bool `db:"real" json:"real,omitempty"`

	// Name either contains the usage locations name or a different descriptor
	Name *pgtype.Text `db:"name" json:"name,omitempty"`

	// LegalPurpose contains the legal purpose for the usage location
	LegalPurpose *[]string `db:"legal_purpose" json:"legalPurpose,omitempty"`

	// MapExcerpt contains the identification of the area the location is in on
	// a topological map using a 1:25000 scale
	MapExcerpt *NumericKeyedValue `db:"map_excerpt" json:"mapExcerpt,omitempty"`

	// MunicipalArea contains the ARS of the municipal and the name in which
	// this usage location is located
	MunicipalArea *NumericKeyedValue `db:"municipal_area" json:"municipalArea,omitempty"`

	// County contains the name of the county the usage location is located in
	County *pgtype.Text `db:"county" json:"county,omitempty"`

	// LandRecord contains information about the parish the usage location is
	// located in
	LandRecord *LandRecord `db:"land_record" json:"landRecord,omitempty"`

	// Plot contains information about the plot the usage location is located
	// in
	Plot *pgtype.Text `db:"plot" json:"plot,omitempty"`

	// MaintenanceAssociation contains information about the (legal) person
	// responsible for maintaining the usage location
	MaintenanceAssociation *NumericKeyedValue `db:"maintenance_association" json:"maintenanceAssociation,omitempty"`

	// EUSurveyArea contains information about an EU-specified survey area
	EUSurveyArea *NumericKeyedValue `db:"eu_survey_area" json:"euSurveyArea,omitempty"`

	// CatchmentAreaCode contains further information about the area the usage
	// location is located in
	CatchmentAreaCode *NumericKeyedValue `db:"catchment_area_code" json:"catchmentAreaCode,omitempty"`

	// RegulationCitation contains a citation from the regulations about water
	// rights
	RegulationCitation *pgtype.Text `db:"regulation_citation" json:"regulationCitation,omitempty"`

	// WithdrawalRates contains information about the rates of water withdrawal
	// that are reflected by the associated water right
	WithdrawalRates []Rate `db:"withdrawal_rates" json:"withdrawalRates,omitempty"`

	// PumpingRates contains information about the allowed pumping rates for the
	// usage location
	PumpingRates []Rate `db:"pumping_rates" json:"pumpingRates,omitempty"`

	// InjectionRates contains information about the allowed rates for
	// injecting something at the usage location
	InjectionRates []Rate `db:"injection_rates" json:"injectionRates,omitempty"`

	// WasteWaterFlowVolume contains information about the amount of waste water
	// handled at that location
	WasteWaterFlowVolume []Rate `db:"waste_water_flow_volume" json:"wasteWaterFlowVolume,omitempty"`

	// RiverBasin describes the river used in connection with the location
	RiverBasin *pgtype.Text `db:"river_basin" json:"riverBasin,omitempty"`

	// GroundwaterBody describes the groundwater used at this location
	GroundwaterBody *pgtype.Text `db:"groundwater_body" json:"groundwaterBody,omitempty"`

	// WaterBody describes the water body used at this location
	WaterBody *pgtype.Text `db:"water_body" json:"waterBody,omitempty"`

	// FloodArea describes the area in which a flood may happen
	FloodArea *pgtype.Text `db:"flood_area" json:"floodArea,omitempty"`

	// WaterProtectionArea contains information about a possible water
	// protection area
	WaterProtectionArea *pgtype.Text `db:"water_protection_area" json:"waterProtectionArea,omitempty"`

	// DamTargetLevels contains information about the levels that are to be
	// expected at a dam
	DamTargetLevels *DamTarget `db:"dam_target_levels" json:"damTargetLevels,omitempty"`

	// FluidDischarge contains rates of discharged fluids
	FluidDischarge []Rate `db:"fluid_discharge" json:"fluidDischarge,omitempty"`

	// RainSupplement contains information about additionally used rain water
	RainSupplement []Rate `db:"rain_supplement" json:"rainSupplement,omitempty"`

	// IrrigationArea contains information about the area that the rain is
	// collected from
	IrrigationArea *Quantity `db:"irrigation_area" json:"irrigationArea,omitempty"`

	// PHValues contains information about the ph values at the usage location
	PHValues *pgtype.Range[pgtype.Numeric] `db:"ph_values" json:"phValues,omitempty"`

	// InjectionLimits contains information about injection limitations
	InjectionLimits []InjectionLimit `db:"injection_limits" json:"injectionLimits,omitempty"`

	// Location contains the GeoJSON representation of the usage locations
	// location
	Location json.RawMessage `db:"location" json:"location"`
}

func (l UsageLocation) MarshalJSON() ([]byte, error) {

	out := usageLocation{
		ID:                     l.ID,
		LocationNumber:         l.LocationNumber,
		SerialID:               l.SerialID,
		WaterRightID:           l.WaterRightID,
		LegalDepartment:        l.LegalDepartment,
		Active:                 l.Active,
		Real:                   l.Real,
		Name:                   l.Name,
		LegalPurpose:           l.LegalPurpose,
		MapExcerpt:             l.MapExcerpt,
		MunicipalArea:          l.MunicipalArea,
		County:                 l.County,
		LandRecord:             l.LandRecord,
		Plot:                   l.Plot,
		MaintenanceAssociation: l.MaintenanceAssociation,
		EUSurveyArea:           l.EUSurveyArea,
		CatchmentAreaCode:      l.CatchmentAreaCode,
		RegulationCitation:     l.RegulationCitation,
		WithdrawalRates:        l.WithdrawalRates,
		PumpingRates:           l.PumpingRates,
		InjectionRates:         l.InjectionRates,
		WasteWaterFlowVolume:   l.WasteWaterFlowVolume,
		RiverBasin:             l.RiverBasin,
		GroundwaterBody:        l.GroundwaterBody,
		WaterBody:              l.WaterBody,
		FloodArea:              l.FloodArea,
		WaterProtectionArea:    l.WaterProtectionArea,
		DamTargetLevels:        l.DamTargetLevels,
		FluidDischarge:         l.FluidDischarge,
		RainSupplement:         l.RainSupplement,
		IrrigationArea:         l.IrrigationArea,
		PHValues:               l.PHValues,
		InjectionLimits:        l.InjectionLimits,
	}

	geom.TransformInPlace(l.Location, func(c geom.Coord) {
		transformer := wgs84.Transform(wgs84.EPSG(25832), wgs84.EPSG(4326))

		lat, long, _ := transformer(c.X(), c.Y(), 0)

		newCoords := geom.Coord{lat, long}
		c.Set(newCoords)
	})

	out.Location, _ = geojson.Marshal(l.Location, geojson.EncodeGeometryWithMaxDecimalDigits(15))
	return json.Marshal(out)
}
