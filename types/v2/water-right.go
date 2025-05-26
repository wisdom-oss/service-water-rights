package v2

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type WaterRight struct {
	Identifiers struct {
		Database uint64  `db:"id"                  json:"database"`
		Cadenza  uint64  `db:"water_right_number"  json:"cadenza"`
		External *string `db:"external_identifier" json:"external"`
		File     *string `db:"file_reference"      json:"fileReference"`
	} `db:"" json:"identifiers"`

	LegalTitle       *string      `db:"legal_title"       json:"legalTitle"`
	Holder           *string      `db:"holder"            json:"holder"`
	Status           *string      `db:"status"            json:"status"`
	InitiallyGranted *pgtype.Date `db:"initially_granted" json:"initiallyGranted"`
	LastChange       *pgtype.Date `db:"last_change"       json:"lastChange"`
	Subject          *string      `db:"subject"           json:"subject"`
	Address          *string      `db:"address"           json:"address"`
	LegalDepartments []string     `db:"legal_departments" json:"legalDepartments"`
	Annotation       *string      `db:"annotation"        json:"annotation"`

	Authorities struct {
		Water       *string `db:"water_authority"       json:"water"`
		Registering *string `db:"registering_authority" json:"registering"`
		Granting    *string `db:"granting_authority"    json:"granting"`
	} `db:"" json:"authorities"`

	Validity struct {
		From  pgtype.Date `db:"valid_from"  json:"from"`
		Until pgtype.Date `db:"valid_until" json:"until"`
	} `db:"" json:"valid"`

	AssociatedUsageLocations []UsageLocation `db:"-" json:"-"`
}

type waterRight struct {
	Identifiers struct {
		Database uint64  `db:"id"                  json:"database"`
		Cadenza  uint64  `db:"water_right_number"  json:"cadenza"`
		External *string `db:"external_identifier" json:"external"`
		File     *string `db:"file_reference"      json:"fileReference"`
	} `db:"" json:"identifiers"`

	LegalTitle       *string      `db:"legal_title"       json:"legalTitle"`
	Holder           *string      `db:"holder"            json:"holder"`
	Status           *string      `db:"status"            json:"status"`
	InitiallyGranted *pgtype.Date `db:"initially_granted" json:"initiallyGranted"`
	LastChange       *pgtype.Date `db:"last_change"       json:"lastChange"`
	Subject          *string      `db:"subject"           json:"subject"`
	Address          *string      `db:"address"           json:"address"`
	LegalDepartments []string     `db:"legal_departments" json:"legalDepartments"`
	Annotation       *string      `db:"annotation"        json:"annotation"`

	Authorities struct {
		Water       *string `db:"water_authority"       json:"water"`
		Registering *string `db:"registering_authority" json:"registering"`
		Granting    *string `db:"granting_authority"    json:"granting"`
	} `db:"" json:"authorities"`

	Validity struct {
		From  pgtype.Date `db:"valid_from"  json:"from"`
		Until pgtype.Date `db:"valid_until" json:"until"`
	} `db:"" json:"valid"`

	AssociatedUsageLocations json.RawMessage `db:"-" json:"usageLocations"`
}

func (r WaterRight) MarshalJSON() ([]byte, error) {
	out := waterRight{
		Identifiers:      r.Identifiers,
		LegalTitle:       r.LegalTitle,
		Holder:           r.Holder,
		Status:           r.Status,
		InitiallyGranted: r.InitiallyGranted,
		LastChange:       r.LastChange,
		Subject:          r.Subject,
		Address:          r.Address,
		LegalDepartments: r.LegalDepartments,
		Annotation:       r.Annotation,
		Authorities:      r.Authorities,
		Validity:         r.Validity,
	}

	featureCollection := geojson.FeatureCollection{
		BBox:     geom.NewBounds(geom.XY),
		Features: make([]*geojson.Feature, 0),
	}

	for _, location := range r.AssociatedUsageLocations {
		geometry := location.EPSG4326Geom()
		featureCollection.BBox.Extend(geometry)

		feature, _ := location.ToFeature()
		featureCollection.Features = append(featureCollection.Features, feature)
	}

	out.AssociatedUsageLocations, _ = featureCollection.MarshalJSON()

	return json.Marshal(out)
}
