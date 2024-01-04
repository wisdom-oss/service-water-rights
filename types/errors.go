package types

import "errors"

var ErrUnsupportedScanInput = errors.New("unsupported scan input. only []byte and string are supported")
var ErrRegexCompileFailed = errors.New("unable to compile regular expression for scan")
var ErrMatchCountWrong = errors.New("unexpected number of regex matches")
