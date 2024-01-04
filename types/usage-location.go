package types

// UsageLocation contains the UsageLocationBaseData as well as further
// information about a single usage location.
type UsageLocation struct {
	// Use the already defined UsageLocationBaseData
	UsageLocationBaseData

	// BasinNumber contains the identification number of the basin used to
	// access the water
	BasinNumber *NumericKeyedName `json:"basinNo,omitempty" db:"basin_no"`

	// County contains the name of the county the usage location is placed in
	County *string `json:"county,omitempty" db:"county"`

	// EUSurveyArea contains the NUTS key used in EU surveys to identify a
	// region
	EUSurveyArea *NumericKeyedName `json:"euSurveyArea,omitempty" db:"eu_survey_area"`

	// Field contains an identification about the field the usage location is
	// placed in
	Field *string `json:"field,omitempty" db:"field"`

	// TODO: Extend description
	GroundwaterVolume *string `json:"groundwaterVolume,omitempty" db:"groundwater_volume"`

	// TODO: Extend description
	LegalScope *string `json:"legalScope,omitempty" db:"legal_scope"`

	// LocalSubDistrict contains the County's district name the
	// usage location is placed in
	LocalSubDistrict *string `json:"localSubDistrict,omitempty" db:"local_sub_district"`

	// MaintenanceAssociation contains the association responsible for
	// maintaining the usage location
	MaintenanceAssociation *NumericKeyedName `json:"maintenanceAssociation,omitempty" db:"maintenance_association"`

	// MunicipalArea contains information about the municipal the usage location
	// is places in
	MunicipalArea *NumericKeyedName `json:"municipalArea,omitempty" db:"municipal_area"`

	// TODO: Extend description
	Plot *string `json:"plot,omitempty" db:"plot"`

	// TODO: Extend description
	Rivershed *string `json:"rivershed,omitempty" db:"rivershed"`

	// SerialNumber contains a string identifying the usage location
	SerialNumber *string `json:"serialNumber,omitempty" db:"serial_no"`

	// TODO: Extend description
	TopMap1To25000 *NumericKeyedName `json:"topMap1To25000,omitempty" db:"top_map_1_25000"`

	// TODO: Extend description
	WaterBody *string `json:"waterBody,omitempty" db:"water_body"`

	// TODO: Extend description
	FloodArea *string `json:"floodArea,omitempty" db:"flood_area"`

	// TODO: Extend description
	WaterProtectionArea *string `json:"waterProtectionArea,omitempty" db:"water_protection_area"`

	// TODO: Extend description
	WithdrawalRate *IntervalRates `json:"withdrawalRates,omitempty" db:"withdrawal_rate"`

	// TODO: Extend description
	FluidDischarge *IntervalRates `json:"fluidDischargeRates,omitempty" db:"fluid_discharge"`

	// TODO: Extend description
	IrrigationArea *IntervalRates `json:"irrigationArea,omitempty" db:"irrigation_area"`

	// TODO: Extend description
	RainSupplement *IntervalRates `json:"rainSupplement,omitempty" db:"rain_supplement"`
}
