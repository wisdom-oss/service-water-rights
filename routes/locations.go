package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/blockloop/scan/v2"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v3"

	"microservice/globals"
	"microservice/types"
)

const (
	locationFilter uint8 = 1 << iota
	stateFilter
	realityFilter
)

// UsageLocations returns a possibly filtered list of usage locations
func UsageLocations(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(wisdomMiddleware.ERROR_CHANNEL_NAME).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.STATUS_CHANNEL_NAME).(<-chan bool)

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

	// now build an array of the arguments and build the query
	var arguments []interface{}
	queryString, err := globals.SqlQueries.Raw("usage-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to load base query: %w", err)
		<-statusChannel
		return
	}

	if enabledFilters&locationFilter != 0 {
		filter, err := globals.SqlQueries.Raw("filter-locations")
		if err != nil {
			errorHandler <- fmt.Errorf("unable to load filter query: %w", err)
			<-statusChannel
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf("WHERE %s", filter)
		} else {
			filter = fmt.Sprintf("AND %s", filter)
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
			<-statusChannel
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf("WHERE %s", filter)
		} else {
			filter = fmt.Sprintf("AND %s", filter)
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
			<-statusChannel
			return
		}

		if len(arguments) == 0 {
			filter = fmt.Sprintf("WHERE %s", filter)
		} else {
			filter = fmt.Sprintf("AND %s", filter)
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

	// now prepare the query
	query, err := globals.Db.Prepare(queryString)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to preparse query: %w", err)
		<-statusChannel
		return
	}

	var rows *sql.Rows
	// now query the database with the correct number of arguments
	if len(arguments) == 0 {
		rows, err = query.Query()
	} else {
		rows, err = query.Query(arguments...)
	}

	if err != nil {
		errorHandler <- fmt.Errorf("error while querying the database: %w", err)
		<-statusChannel
		return
	}

	var locations []types.UsageLocation
	err = scan.Rows(&locations, rows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(204)
			return
		}
		errorHandler <- fmt.Errorf("unable to parse query result: %w", err)
		<-statusChannel
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to return response: %w", err)
		<-statusChannel
		return
	}
}
