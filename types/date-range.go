package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const dateRangeRegEx = `^\[(\d{4}-\d{2}-\d{2}|[+-]?infinity),(\d{4}-\d{2}-\d{2}|[+-]?infinity)?\)$`

var fromInfinityDate = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
var untilInfinityDate = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

type DateRange struct {
	From  time.Time `json:"from"`
	Until time.Time `json:"until"`
}

func (dr *DateRange) Scan(src interface{}) error {
	var value string
	switch src.(type) {
	case []byte:
		value = string(src.([]byte))
		break
	case string:
		value = src.(string)
		break
	default:
		return ErrUnsupportedScanInput
	}
	if value == "empty" {
		return nil
	}
	regex, err := regexp.Compile(dateRangeRegEx)
	if err != nil {
		err := errors.Wrap(err, ErrRegexCompileFailed.Error())
		return err
	}
	matches := regex.FindStringSubmatch(value)
	switch len(matches) {
	case 2:
		matches = append(matches, "infinity")
		break
	case 3:
		if matches[2] == "" {
			matches[2] = "infinity"
		}
	default:
		return fmt.Errorf("%w: expected two (2) or three (3) matches, got %d", ErrMatchCountWrong, len(matches))
	}
	for idx, match := range matches {
		matches[idx] = strings.TrimSpace(match)
	}

	if matches[1] == "infinity" {
		dr.From = fromInfinityDate
	} else {
		dr.From, err = time.Parse("2006-01-02", matches[1])
		if err != nil {
			return err
		}
	}

	if matches[2] == "infinity" {
		dr.Until = untilInfinityDate
	} else {
		dr.Until, err = time.Parse("2006-01-02", matches[2])
		if err != nil {
			return err
		}
	}
	return nil
}
