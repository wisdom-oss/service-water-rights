package v2

import (
	"encoding/json"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/wroge/wgs84/v2"
)

// UsageLocation represents a location that has been crawled from the Cadenza
// Database.
type UsageLocation struct {
	ID                  int                   `db:"id"                      json:"internalID"`
	CadenzaID           int                   `db:"no"                      json:"cadenzaID"`
	WaterRightID        int                   `db:"water_right"             json:"waterRightID"`
	Serial              *string               `db:"serial"                  json:"serial"`
	Active              *bool                 `db:"active"                  json:"isActive"`
	Virtual             *bool                 `db:"real"                    json:"isVirtual"`
	Name                *string               `db:"name"                    json:"name"`
	LegalDepartment     *string               `db:"legal_department"        json:"legalDepartment"`
	LegalPurpose        *[]string             `db:"legal_purpose"           json:"legalPurposes"`
	MapExcerpt          *NumericKeyedValue    `db:"map_excerpt"             json:"mapExcerpt"`
	MunicipalArea       *NumericKeyedValue    `db:"municipal_area"          json:"municipalArea"`
	County              *string               `db:"county"                  json:"county"`
	Plot                *string               `db:"plot"                    json:"plot"`
	Maintenance         *NumericKeyedValue    `db:"maintenance_association" json:"maintenance"`
	SurveyArea          *NumericKeyedValue    `db:"eu_survey_area"          json:"surveyArea"`
	CatchmentArea       *NumericKeyedValue    `db:"catchment_area_code"     json:"catchmentArea"`
	RegulationCitation  *string               `db:"regulation_citation"     json:"regulation"`
	GroundwaterBody     *string               `db:"groundwater_body"        json:"groundwaterBody"`
	WaterBody           *string               `db:"water_body"              json:"waterBody"`
	FloodArea           *string               `db:"flood_area"              json:"floodArea"`
	WaterProtectionArea *string               `db:"water_protection_area"   json:"waterProtectionArea"`
	RiverBasin          *string               `db:"river_basin"             json:"riverBasin"`
	PhValues            pgtype.Range[float64] `db:"ph_values"               json:"phValues"`
	InjectionLimits     []InjectionLimit      `db:"injection_limits"        json:"injectionLimits"`
	LandRecord          *LandRecord           `db:"land_record"             json:"landRecord"`
	IrrigationArea      *Quantity             `db:"irrigation_area"         json:"irrigationArea"`
	DamTargetLevels     *DamTarget            `db:"dam_target_levels"       json:"damTargetLevels"`
	Rates               struct {
		Withdrawal      []Rate `db:"withdrawal_rates"        json:"withdrawal"`
		Pumping         []Rate `db:"pumping_rates"           json:"pumping"`
		Injection       []Rate `db:"injection_rates"         json:"injection"`
		WasteWater      []Rate `db:"waste_water_flow_volume" json:"wasteWaterFlow"`
		FluidDischarges []Rate `db:"fluid_discharge"         json:"fluidDischarges"`
		RainSupplements []Rate `db:"rain_supplement"         json:"rainSupplement"`
	} `db:"" json:"rates"`
	Geometry geom.T `db:"location" json:"-"`
}

func (l UsageLocation) EPSG4326Geom() geom.T {
	if l.Geometry.SRID() == defaultCRS {
		return l.Geometry
	}
	geom.TransformInPlace(l.Geometry, func(c geom.Coord) {
		transformer := wgs84.Transform(wgs84.EPSG(l.Geometry.SRID()), wgs84.EPSG(defaultCRS))

		lat, long, _ := transformer(c.X(), c.Y(), 0)

		newCoords := geom.Coord{lat, long}
		c.Set(newCoords)
	})
	l.Geometry, _ = geom.SetSRID(l.Geometry, defaultCRS)
	return l.Geometry
}

func (l UsageLocation) ToFeature() (*geojson.Feature, error) {
	if l.Geometry.SRID() != defaultCRS {
		l.EPSG4326Geom()
	}

	feature := &geojson.Feature{
		ID:       strconv.Itoa(l.ID),
		Geometry: l.Geometry,
	}

	marshalled, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	var properties map[string]any
	if err := json.Unmarshal(marshalled, &properties); err != nil {
		return nil, err
	}

	feature.Properties = properties
	feature.Properties["id"] = strconv.Itoa(l.ID)
	if feature.Properties["virtual"] != nil {
		feature.Properties["virtual"] = !*l.Virtual
	}
	return feature, nil
}
