package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/blockloop/scan/v2"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v3"

	"microservice/globals"
	"microservice/types"
)

// WaterRights returns all water rights stored in the database with all usage
// locations.
func WaterRights(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(wisdomMiddleware.ERROR_CHANNEL_NAME).(chan<- interface{})
	statusChannel := r.Context().Value(wisdomMiddleware.STATUS_CHANNEL_NAME).(<-chan bool)

	// now query the database for all water rights available in the database
	rows, err := globals.SqlQueries.Query(globals.Db, "water-rights")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to get all water rights: %w", err)
		<-statusChannel
		return
	}
	var waterRights []types.WaterRight
	err = scan.Rows(&waterRights, rows)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to parse database rows: %w", err)
		<-statusChannel
		return
	}

	// now build the filter for the usage location request for the water rights
	filterQuery := `WHERE water_right = $1::int;`
	baseQuery, err := globals.SqlQueries.Raw("extended-usage-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to load query: %w", err)
		<-statusChannel
		return
	}
	queryString := strings.ReplaceAll(baseQuery, `;`, "")
	queryString += fmt.Sprintf(" %s", filterQuery)

	// now prepare the query
	query, err := globals.Db.Prepare(queryString)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to prepare sql query: %w", err)
		<-statusChannel
		return
	}

	for idx, waterRight := range waterRights {
		// execute the prepared query
		rows, err := query.Query(waterRight.NlwknId)
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
		waterRights[idx] = waterRight
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(waterRights)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to send response: %w", err)
		<-statusChannel
		return
	}
}
