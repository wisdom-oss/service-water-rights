package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v3"

	"microservice/globals"
	"microservice/types"
)

// SingleWaterRight returns the water right referenced by the id supplied in the
// request url
func SingleWaterRight(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(wisdomMiddleware.ERROR_CHANNEL_NAME).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.STATUS_CHANNEL_NAME).(<-chan bool)

	// get the water right number from the url
	waterRightNumber := chi.URLParam(r, "water-right-nlwkn-id")

	// now build the query that filters the water right to only contain the
	// specified water right
	filterQuery := `WHERE no = $1::int;`
	baseQuery, err := globals.SqlQueries.Raw("water-rights")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to load base query for water rights: %w", err)
		<-statusChannel
		return
	}
	queryString := strings.ReplaceAll(baseQuery, `;`, "")
	queryString += fmt.Sprintf(" %s", filterQuery)
	query, err := globals.Db.Prepare(queryString)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to prepare sql query: %w", err)
		<-statusChannel
		return
	}

	rows, err := query.Query(waterRightNumber)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to get all water rights: %w", err)
		<-statusChannel
		return
	}

	var waterRight types.WaterRight
	err = scan.Row(&waterRight, rows)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler <- "UNKNOWN_NLWKN_ID"
			<-statusChannel
			return
		}
		errorHandler <- fmt.Errorf("unable to parse database rows: %w", err)
		<-statusChannel
		return
	}

	// now build the filter for the usage location request for the water rights
	filterQuery = `WHERE water_right = $1::int;`
	baseQuery, err = globals.SqlQueries.Raw("extended-usage-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to load query: %w", err)
		<-statusChannel
		return
	}
	queryString = strings.ReplaceAll(baseQuery, `;`, "")
	queryString += fmt.Sprintf(" %s", filterQuery)

	// now prepare the query
	query, err = globals.Db.Prepare(queryString)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to prepare sql query: %w", err)
		<-statusChannel
		return
	}

	// execute the prepared query
	rows, err = query.Query(waterRight.NlwknId)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to query usage locations: %w", err)
		<-statusChannel
		return
	}
	// now parse the usage locations
	var usageLocations []types.UsageLocation
	err = scan.RowsStrict(&usageLocations, rows)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to parse database rows: %w", err)
		<-statusChannel
		return
	}

	// now set the parsed usage locations on the water right
	waterRight.UsageLocations = usageLocations

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(waterRight)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to send response: %w", err)
		<-statusChannel
		return
	}
}
