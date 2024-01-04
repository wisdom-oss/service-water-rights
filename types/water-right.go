package types

import "time"

// WaterRight collects all UsageLocations of the single water right and contains
// legal information about the wate right and the last modification date, if
// available
type WaterRight struct {
	// ID contains the internally used id to reference this water right
	ID int `json:"id" db:"id"`

	// NlwknId contains the id used by the NLWKN to reference the water right
	NlwknId int `json:"nlwknID" db:"no"`

	// ExternalID contains the id that is used by the NLWKN to reference the
	// water right in external applications
	ExternalID *string `json:"externalID,omitempty" db:"ext_id"`

	// FileReference contains a reference to the file containing the water
	// right
	FileReference *string `json:"fileReference,omitempty" db:"file_ref"`

	// State contains the current state of the water right
	State *string `json:"state,omitempty" db:"state"`

	// Subject contains a subject listed for the water right
	Subject *string `json:"subject,omitempty" db:"subject"`

	// Address contains the user's address
	Address *string `json:"address,omitempty" db:"address"`

	// Annotation contains an annotation related to the water right which may
	// further limit the rights granted
	Annotation *string `json:"annotation,omitempty" db:"address"`

	// Bailee contains the information about the keeper of the water right
	Bailee *string `json:"bailee,omitempty" db:"bailee"`

	// LastChange contains the information about the last time the water right
	// has been changed
	LastChange *time.Time `json:"dateOfChange,omitempty" db:"date_of_change"`

	// Validity contains a DateRange showing in which time range the water right
	// is valid
	Validity *DateRange `json:"valid,omitempty" db:"valid"`

	// GrantingAuthority contains the name of the authority that granted the
	// water right
	GrantingAuthority *string `json:"grantingAuthority,omitempty" db:"granting_authority"`

	// RegisteringAuthority contains the name of the authority that the water
	// right has been registered with
	RegisteringAuthority *string `json:"registeringAuthority,omitempty" db:"registering_authority"`

	// WaterAuthority contains the name of the water authority related to the
	// water right
	WaterAuthority *string `json:"waterAuthority,omitempty" db:"water_authority"`

	// UsageLocations contains all usage locations that are listed for this
	// water right
	UsageLocations []UsageLocation `json:"locations" db:"-"`
}
