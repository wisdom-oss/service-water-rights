package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const numericKeyedNameRegEx = `^\((\d+),"([^{}[\]]+)"\)$`

type NumericKeyedName struct {
	Key  int64  `json:"key"`
	Name string `json:"name"`
}

func (nkn *NumericKeyedName) Scan(src interface{}) error {
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
	regex, err := regexp.Compile(numericKeyedNameRegEx)
	if err != nil {
		err := errors.Wrap(err, ErrRegexCompileFailed.Error())
		return err
	}
	matches := regex.FindStringSubmatch(value)
	if len(matches) != 3 {
		return fmt.Errorf("%w: expected 3 matches, got %d", ErrMatchCountWrong, len(matches))
	}
	nkn.Key, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return err
	}
	nkn.Name = strings.TrimSpace(matches[2])
	return nil
}
