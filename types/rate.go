package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const rateRegEx = `^\(([+-]?\d+\.?\d?),"([\w\s²³'"]+)"\)$`

// Rate describes the amount of a medium that flows in the span of a time frame
// denoted in the unit
type Rate struct {
	// Amount contains the numerical part of the rate
	Amount float64 `json:"amount"`
	// Unit contains the unit denoting the time frame and the unit for the
	// amount
	Unit string `json:"unit"`
}

// Scan implements the Scanner interface allowing reading the results returned
// by the database into the data type.
// The Scan function only supports []byte and sting as input since the value
// is extracted using a regular expression
func (r *Rate) Scan(src interface{}) error {
	// check the type of the database output and convert it into a string if
	// it may be converted
	var value string
	switch src.(type) {
	case []byte:
		value = string(src.([]byte))
		break
	case string:
		value = src.(string)
	default:
		return ErrUnsupportedScanInput
	}
	regex, err := regexp.Compile(rateRegEx)
	if err != nil {
		err := errors.Wrap(err, ErrRegexCompileFailed.Error())
		return err
	}
	matches := regex.FindStringSubmatch(value)
	// now parse the values listed in the matches if the number of matches is
	// correct
	if len(matches) != 3 {
		return fmt.Errorf("%w: expected 3 matches, got %d", ErrMatchCountWrong, len(matches))
	}
	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return fmt.Errorf("unable to parse 'amount' field as float: %w", err)
	}
	r.Amount = amount
	r.Unit = strings.TrimSpace(matches[2])
	return nil
}
