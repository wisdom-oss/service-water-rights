package types

import "strings"

type IntervalRates []IntervalRate

func (irs *IntervalRates) Scan(src interface{}) error {
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

	// now clean up the row to only contain the array values
	value = strings.Trim(value, `{}"`)
	// now split the array values to retrive the single entries
	singleEntries := strings.Split(value, `","`)
	// now iterate over the single entries and read them in
	for _, singleEntry := range singleEntries {
		// cleanup the single entry
		singleEntry = strings.Trim(singleEntry, `"`)
		// now create a interval rate from it
		ir := IntervalRate{}
		err := ir.Scan(singleEntry)
		if err != nil {
			return err
		}
		*irs = append(*irs, ir)
	}
	return nil
}
