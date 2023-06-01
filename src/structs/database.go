package structs

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	geojson "github.com/paulmach/go.geojson"
)

// Rate is the golang implementation of the postgres composite type rate
type Rate struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

func (r *Rate) FromJSON(input pgtype.JSON) error {
	err := json.Unmarshal(input.Bytes, r)
	if err != nil {
		return err
	}
	return nil
}

// IntervalRate is the golang implementation of the postgres composite type interval_rate
type IntervalRate struct {
	Rate
	Duration *string `json:"duration"`
}

func (iR *IntervalRate) FromJSON(input pgtype.JSON) error {
	err := json.Unmarshal(input.Bytes, iR)
	if err != nil {
		return err
	}
	return nil
}

// NumericKeyedName is the golang implementation of the postgres composite type numeric_keyed_name
type NumericKeyedName struct {
	Key  *int    `json:"key"`
	Name *string `json:"name"`
}

func (nkm *NumericKeyedName) FromJSON(input pgtype.JSON) error {
	err := json.Unmarshal(input.Bytes, nkm)
	if err != nil {
		return err
	}
	return nil
}

// DateTimeRange the golang implementation of the postgres datatype Daterange
type DateTimeRange struct {
	From  *int64 `json:"from"`
	Until *int64 `json:"until"`
}

// DbWaterRight represents a single water right stored in the database
type DbWaterRight struct {
	ID                   int               `db:"id"`
	NlwknID              int               `db:"no"`
	ExternalID           *string           `db:"ext_id"`
	FileReference        *string           `db:"file_ref"`
	LegalTitle           *string           `db:"legal_title"`
	State                *string           `db:"state"`
	Subject              *string           `db:"subject"`
	Address              *string           `db:"address"`
	Annotation           *string           `db:"annotation"`
	Bailee               *string           `db:"bailee"`
	DateOfChange         *pgtype.Timestamp `db:"date_of_change"`
	Valid                *pgtype.Daterange `db:"valid"`
	GrantingAuthority    *string           `db:"granting_authority"`
	RegisteringAuthority *string           `db:"registering_authority"`
	WaterAuthority       *string           `db:"water_authority"`
}

// ToWaterRight converts the water right that has been pulled from the database
// into a WaterRight that may be sent as a response for a request of the
// water right. It also handles the conversion of postgres types into golang
// types or structs
func (r DbWaterRight) ToWaterRight() WaterRight {
	// convert the date of change to a unix timestamp
	var dateOfChangeTimestamp *int64
	if r.DateOfChange == nil {
		dateOfChangeTimestamp = nil
	} else {
		timestamp := r.DateOfChange.Time.Unix()
		dateOfChangeTimestamp = &timestamp
	}

	// convert the valid date range into a new object
	var validty *DateTimeRange
	if r.Valid == nil {
		validty = nil
	} else {
		lowerBound := r.Valid.Lower.Time.Unix()
		upperBound := r.Valid.Upper.Time.Unix()
		validty = &DateTimeRange{
			From:  &lowerBound,
			Until: &upperBound,
		}
	}

	return WaterRight{
		ID:                   r.ID,
		NlwknID:              r.NlwknID,
		ExternalId:           r.ExternalID,
		FileReference:        r.FileReference,
		LegalTitle:           r.LegalTitle,
		State:                r.State,
		Subject:              r.Subject,
		Address:              r.Address,
		Annotation:           r.Annotation,
		Bailee:               r.Bailee,
		DateOfChange:         dateOfChangeTimestamp,
		Valid:                validty,
		GrantingAuthority:    r.GrantingAuthority,
		RegisteringAuthority: r.RegisteringAuthority,
		WaterAuthority:       r.WaterAuthority,
	}
}

type DbUsageLocation struct {
	ID                     int               `db:"id"`
	WaterRight             int               `db:"water_right"`
	Name                   *string           `db:"name"`
	No                     *int              `db:"no"`
	Active                 *pgtype.Bool      `db:"active"`
	Location               *geojson.Geometry `db:"location"`
	BasinNo                *pgtype.JSON      `db:"basin_no"`
	County                 *string           `db:"county"`
	EuSurveyArea           *pgtype.JSON      `db:"eu_survey_area"`
	Field                  *int              `db:"field"`
	GroundwaterVolume      *string           `db:"groundwater_volume"`
	LegalScope             *string           `db:"legal_scope"`
	LocalSubDistrict       *string           `db:"local_sub_district"`
	MaintenanceAssociation *pgtype.JSON      `db:"maintenance_association"`
	MunicipalArea          *pgtype.JSON      `db:"municipal_area"`
	Plot                   *string           `db:"plot"`
	Real                   *pgtype.Bool      `db:"real"`
	Rivershed              *string           `db:"rivershed"`
	SerialNo               *string           `db:"serial_no"`
	TopMap1To25000         *pgtype.JSON      `db:"top_map_1_25000"`
	WaterBody              *string           `db:"water_body"`
	FloodArea              *string           `db:"flood_area"`
	WaterProtectionArea    *string           `db:"water_protection_area"`
	WithdrawalRates        *pgtype.JSON      `db:"withdrawal_rate"`
	FluidDischarge         *pgtype.JSON      `db:"fluid_discharge"`
	IrrigationArea         *pgtype.JSON      `db:"irrigation_area"`
	RainSupplement         *pgtype.JSON      `db:"rain_supplement"`
}

// ToUsageLocation converts the database entry for a usage location to a basic
// usage location that may be used when getting usage locations in a broader
// area
func (l DbUsageLocation) ToUsageLocation() UsageLocation {
	// convert the postgres boolean to golang
	var activeLocation *bool
	if l.Active == nil {
		activeLocation = nil
	} else {
		b := l.Active.Bool
		activeLocation = &b
	}

	// convert the postgres boolean to golang
	var realLocation *bool
	if l.Real == nil {
		realLocation = nil
	} else {
		b := l.Real.Bool
		realLocation = &b
	}
	return UsageLocation{
		ID:         l.ID,
		Name:       l.Name,
		WaterRight: l.WaterRight,
		IsActive:   activeLocation,
		IsReal:     realLocation,
		Location:   l.Location,
	}
}

// ToDetailedUsageLocation converts the database entry for a usage location to a
// usage location that may be used when getting usage locations for a single
// water right
func (l DbUsageLocation) ToDetailedUsageLocation() DetailedUsageLocation {
	// convert the basin number to a NumericKeyedName
	var basinNumber *NumericKeyedName
	if l.BasinNo == nil {
		basinNumber = nil
	} else {
		var nkm NumericKeyedName
		err := nkm.FromJSON(*l.BasinNo)
		if err != nil {
			panic(err)
		}
		basinNumber = &nkm
	}

	// convert the eu survey area into a NumericKeyedName
	var euSurveyArea *NumericKeyedName
	if l.EuSurveyArea == nil {
		euSurveyArea = nil
	} else {
		var nkm NumericKeyedName
		err := nkm.FromJSON(*l.EuSurveyArea)
		if err != nil {
			panic(err)
		}
		euSurveyArea = &nkm
	}

	// convert the maintenance association into a NumericKeyedName
	var maintenanceAssoc *NumericKeyedName
	if l.MaintenanceAssociation == nil {
		maintenanceAssoc = nil
	} else {
		var nkm NumericKeyedName
		err := nkm.FromJSON(*l.MaintenanceAssociation)
		if err != nil {
			panic(err)
		}
		maintenanceAssoc = &nkm
	}

	// convert the municipal area into a NumericKeyedName
	var municipalArea *NumericKeyedName
	if l.MunicipalArea == nil {
		municipalArea = nil
	} else {
		var nkm NumericKeyedName
		err := nkm.FromJSON(*l.MunicipalArea)
		if err != nil {
			panic(err)
		}
		municipalArea = &nkm
	}

	// convert the TopMap1To25000 into a NumericKeyedName
	var topMap *NumericKeyedName
	if l.TopMap1To25000 == nil {
		topMap = nil
	} else {
		var nkm NumericKeyedName
		err := nkm.FromJSON(*l.TopMap1To25000)
		if err != nil {
			panic(err)
		}
		topMap = &nkm
	}

	// convert the withdrawal rates
	var withdrawalRates []*IntervalRate
	if l.WithdrawalRates == nil {
		withdrawalRates = nil
	} else {
		err := l.WithdrawalRates.AssignTo(&withdrawalRates)
		if err != nil {
			panic(err)
		}
	}

	// convert the fluid discharge rates
	var fluidDischargeRates []*IntervalRate
	if l.FluidDischarge == nil {
		fluidDischargeRates = nil
	} else {
		err := l.FluidDischarge.AssignTo(&fluidDischargeRates)
		if err != nil {
			panic(err)
		}
	}

	// convert the fluid discharge rates
	var rainSupplementRates []*IntervalRate
	if l.RainSupplement == nil {
		rainSupplementRates = nil
	} else {
		err := l.RainSupplement.AssignTo(&rainSupplementRates)
		if err != nil {
			panic(err)
		}
	}

	// convert the irrigation area rate
	var irrigationAreaRate *Rate
	if l.IrrigationArea == nil {
		irrigationAreaRate = nil
	} else {
		var r Rate
		err := r.FromJSON(*l.IrrigationArea)
		if err != nil {
			panic(err)
		}
		irrigationAreaRate = &r
	}

	// now build the response
	return DetailedUsageLocation{
		ID:                     l.ID,
		Name:                   l.Name,
		WaterRight:             l.WaterRight,
		IsActive:               &l.Active.Bool,
		IsReal:                 &l.Real.Bool,
		Location:               l.Location,
		NlwknID:                l.No,
		BasinNumber:            basinNumber,
		County:                 l.County,
		EuSurveyArea:           euSurveyArea,
		Field:                  l.Field,
		GroundwaterVolume:      l.GroundwaterVolume,
		LegalScope:             l.LegalScope,
		LocalSubDistrict:       l.LocalSubDistrict,
		MaintenanceAssociation: maintenanceAssoc,
		MunicipalArea:          municipalArea,
		Plot:                   l.Plot,
		Rivershed:              l.Rivershed,
		SerialNumber:           l.SerialNo,
		TopMap1To25000:         topMap,
		WaterBody:              l.WaterBody,
		FloodArea:              l.FloodArea,
		WaterProtectionArea:    l.WaterProtectionArea,
		WithdrawalRates:        withdrawalRates,
		FluidDischarge:         fluidDischargeRates,
		IrrigationArea:         irrigationAreaRate,
		RainSupplement:         rainSupplementRates,
	}

}
