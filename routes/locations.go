package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lib/pq"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"github.com/wisdom-oss/service-water-rights/types"

	"github.com/wisdom-oss/service-water-rights/globals"
)

const (
	locationFilter uint8 = 1 << iota
	stateFilter
	realityFilter
)

// UsageLocations returns a possibly filtered list of usage locations
func UsageLocations(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})

	// now check which filters were enabled
	var enabledFilters uint8
	if r.URL.Query().Has("in") {
		enabledFilters = enabledFilters | locationFilter
	}
	if r.URL.Query().Has("is_active") {
		enabledFilters = enabledFilters | stateFilter
	}
	if r.URL.Query().Has("is_real") {
		enabledFilters = enabledFilters | realityFilter
	}

	var arguments []interface{}
	// now build an array of the arguments and build the query
	queryString, err := globals.SqlQueries.Raw("get-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to load base query: %w", err)
		return
	}

	if enabledFilters&locationFilter != 0 {
		filter, err := globals.SqlQueries.Raw("filter-locations")
		if err != nil {
			errorHandler <- fmt.Errorf("unable to load filter query: %w", err)
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf(" WHERE %s", filter)
		} else {
			filter = fmt.Sprintf(" AND %s", filter)
		}

		// now set the correct argument number
		filter = strings.ReplaceAll(filter, "$1", fmt.Sprintf("$%d", len(arguments)+1))
		queryString += filter
		arguments = append(arguments, pq.Array(r.URL.Query()["in"]))
	}

	if enabledFilters&stateFilter != 0 {
		filter, err := globals.SqlQueries.Raw("filter-state")
		if err != nil {
			errorHandler <- fmt.Errorf("unable to load filter query: %w", err)
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf(" WHERE %s", filter)
		} else {
			filter = fmt.Sprintf(" AND %s", filter)
		}

		// now set the correct argument number
		filter = strings.ReplaceAll(filter, "$1", fmt.Sprintf("$%d", len(arguments)+1))
		queryString += filter
		arguments = append(arguments, r.URL.Query()["is_active"][0])
	}

	if enabledFilters&realityFilter != 0 {
		filter, err := globals.SqlQueries.Raw("filter-reality")
		if err != nil {
			errorHandler <- fmt.Errorf("unable to load filter query: %w", err)
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf(" WHERE %s", filter)
		} else {
			filter = fmt.Sprintf(" AND %s", filter)
		}

		// now set the correct argument number
		filter = strings.ReplaceAll(filter, "$1", fmt.Sprintf("$%d", len(arguments)+1))
		queryString += filter
		arguments = append(arguments, r.URL.Query()["is_real"][0])
	}

	// now clean the query string by removing all semicolons and setting one
	// at the end
	queryString = strings.ReplaceAll(queryString, ";", "")
	queryString += ";"

	api, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		errorHandler <- fmt.Errorf("unable to prepare query parser: %w", err)
		return
	}
	scanner, err := pgxscan.NewAPI(api)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to create query parser: %w", err)
	}

	var locations []types.UsageLocation
	err = scanner.Select(r.Context(), globals.Db, &locations, queryString, arguments...)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve usage locations: %w", err)
		return
	}

	if len(locations) == 0 {
		w.WriteHeader(204)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to return response: %w", err)
		return
	}
}
