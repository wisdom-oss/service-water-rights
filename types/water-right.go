package types

import "github.com/jackc/pgx/v5/pgtype"

// WaterRight represents a water right entry, incorporating various details
// about the right including the holder of the right, the time window of its
// validity, the status and the legal title of the right.
// It supplies in-depth specifics such as the granting and registering
// authorities, the first granted and last changed dates of the right, and a
// reference to the water right application.
// It also includes the subject of the right, the address of the right holder,
// any legal departments associated with the right, and other annotations.
type WaterRight struct {
	// ID represents the ID issued for this water right by the database
	ID pgtype.Int8 `db:"id" json:"id"`

	// WaterRightNumber represents the ID of the water right issued by the NLWKN
	WaterRightNumber pgtype.Int8 `db:"water_right_number" json:"water_right_number"`

	// Holder contains the holder's name for this water right
	Holder *pgtype.Text `db:"holder" json:"holder"`

	// ValidFrom contains the date from which on the water right is valid
	ValidFrom *pgtype.Date `db:"valid_from" json:"validFrom"`

	// ValidUntil contains the date until the water right is valid and may be
	// used
	ValidUntil *pgtype.Date `db:"valid_until" json:"validUntil"`

	// Status contains a textual description of the water rights state
	Status *pgtype.Text `db:"status" json:"status"`

	// LegalTitle contains information about the title issued for the water
	// right
	LegalTitle *pgtype.Text `db:"legal_title" json:"legalTitle"`

	// WaterAuthority contains the name of the water authority responsible for
	// the water right
	WaterAuthority *pgtype.Text `db:"water_authority" json:"waterAuthority"`

	// RegisteringAuthority contains the name of the authority that the water
	// right has been registered with
	RegisteringAuthority *pgtype.Text `db:"registering_authority" json:"registeringAuthority"`

	// GrantingAuthority contains the name of the authority that granted the
	// water right
	GrantingAuthority *pgtype.Text `db:"granting_authority" json:"grantingAuthority"`

	// InitiallyGranted contains the date at which the water right has been granted
	// for the first time
	InitiallyGranted *pgtype.Date `db:"initially_granted" json:"initiallyGranted"`

	// LastChange contains the date at which the water right has been changed
	// for the last time
	LastChange *pgtype.Date `db:"last_change" json:"lastChange"`

	// FileReference contains the reference to the water right application
	FileReference *pgtype.Text `db:"file_reference" json:"fileReference"`

	// ExternalIdentifier contains an external identifier assigned by the
	// RegisteringAuthority
	ExternalIdentifier *pgtype.Text `db:"external_identifier" json:"externalIdentifier"`

	// Subject contains the subject of the water right
	Subject *pgtype.Text `db:"subject" json:"subject"`

	// Address contains the address of the RightsHolder
	Address *pgtype.Text `db:"address" json:"address"`

	// LegalDepartments contains the identifiers for the legal departments the
	// water right has been assigned to.
	//
	// The possible values are:
	//   * A: Withdrawal of water or solid substances from surface waters
	//   * B: Introduction and discharge of substances into surface and coastal waters
	//   * C: Damming and lowering of surface waters
	//   * D: Other impact on surface waters
	//   * E: Withdrawal, pumping and discharge of groundwater
	//   * F: Other uses and impacts on groundwater
	//   * K: Compulsory rights
	//   * L: Fishing Rights
	LegalDepartments []string `db:"legal_departments" json:"legalDepartments"`

	// Annotation contains other annotations for the water right
	Annotation *pgtype.Text `db:"annotation" json:"annotation"`
}
