package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sanyokbig/pqinterval"
)

const intervalRateRegEx = `\(([+-]?\d+\.?\d?),\\"([\w\s²³'"]+)\\",(?:\\")?(\d{2}:\d{2}:\d{2}|[\d\s\w]+)(?:\\")?\)`

// IntervalRate reflects a rate that may only be used Every duration.
type IntervalRate struct {
	Rate
	Every    time.Duration       `json:"interval"`
	Interval pqinterval.Interval `json:"-"`
}

func (ir *IntervalRate) Scan(src interface{}) error {
	var value string
	switch src.(type) {
	case []byte:
		value = string(src.([]byte))
	case string:
		value = src.(string)
	default:
		return ErrUnsupportedScanInput
	}
	regex, err := regexp.Compile(intervalRateRegEx)
	if err != nil {
		err := errors.Wrap(err, ErrRegexCompileFailed.Error())
		return err
	}
	matches := regex.FindStringSubmatch(value)
	// now parse the values listed in the matches if the number of matches is
	// correct
	if len(matches) != 4 {
		return fmt.Errorf("%w: expected 4 matches, got %d", ErrMatchCountWrong, len(matches))
	}
	amount, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return fmt.Errorf("unable to parse 'amount' field as float: %w", err)
	}
	ir.Amount = amount
	ir.Unit = strings.TrimSpace(matches[2])
	durationString := strings.ReplaceAll(matches[3], `\"`, "")
	ival := pqinterval.Interval{}
	err = ival.Scan(durationString)
	if err != nil {
		return fmt.Errorf("unable to create entry since the interval could not be parsed correctly: %w", err)
	}
	ir.Interval = ival
	dur, err := ival.Duration()
	if err != nil {
		return fmt.Errorf("unable to create entry since the duration could not be converted to duration correctly"+
			": %w", err)
	}
	ir.Every = dur
	return nil
}

func (ir IntervalRate) MarshalJSON() ([]byte, error) {
	type Output struct {
		Rate
		Interval string `json:"interval"`
	}
	var val interface{}
	strc := Output{Rate: ir.Rate}
	val, _ = ir.Interval.Value()
	interval := val.(string)
	strc.Interval = interval
	return json.Marshal(strc)
}
