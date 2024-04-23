package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/goccy/go-json"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"github.com/wisdom-oss/service-water-rights/globals"
	"github.com/wisdom-oss/service-water-rights/types"
)

// WaterRights returns all water rights stored in the database with all usage
// locations.
func WaterRights(w http.ResponseWriter, r *http.Request) {
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})

	query, err := globals.SqlQueries.Raw("water-rights")
	if err != nil {
		errorHandler <- fmt.Errorf("failed to retrieve query for water rights: %w", err)
		return
	}

	var waterRights []types.WaterRight
	err = pgxscan.Select(r.Context(), globals.Db, &waterRights, query)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to query water rights: %w", err)
		return
	}

	query, err = globals.SqlQueries.Raw("get-water-right-usage-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to query water right locations: %w", err)
		return
	}

	var responseContents []struct {
		types.WaterRight
		Locations []types.UsageLocation `json:"locations"`
	}

	var errors []error
	var successes []bool

	for _, waterRight := range waterRights {
		go func(wr types.WaterRight, query string) {
			var usageLocations []types.UsageLocation
			err = pgxscan.Select(r.Context(), globals.Db, &usageLocations, query)
			if err != nil {
				errors = append(errors, err)
				return
			}
			responseContents = append(responseContents, struct {
				types.WaterRight
				Locations []types.UsageLocation `json:"locations"`
			}{
				WaterRight: waterRight,
				Locations:  usageLocations,
			})
			successes = append(successes, true)

		}(waterRight, query)
	}

	for {
		if len(errors)+len(successes) == len(waterRights) {
			break
		} else {
			time.Sleep(150 * time.Millisecond)
		}
	}

	if len(errors) != 0 {
		for _, err := range errors {
			errorHandler <- fmt.Errorf("unable to query water right locations: %w", err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(responseContents)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to encode response: %w", err)
		return
	}

}
