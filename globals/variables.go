package globals

import (
	"github.com/qustavo/dotsql"
)

// This file contains globally shared variables (e.g., service name, sql queries)

// ServiceName contains the global identifier for the service.
const ServiceName = "water-rights"

// SqlQueries contains the prepared sql queries from the resources folder.
var SqlQueries *dotsql.DotSql

// variables.
var Environment map[string]string = make(map[string]string)
