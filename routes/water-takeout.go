package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"github.com/wisdom-oss/service-water-rights/globals"
	"github.com/wisdom-oss/service-water-rights/types"
)

var ErrNoWithdrawalRatesAvailable = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.6.1",
	Status: 500,
	Title:  "No Withdrawal Rates Available",
	Detail: "In the selected area are no usage locations with a withdrawal rate",
}

func WaterTakeout(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})

	// now parse the json request body
	var inputGeometries []geojson.Geometry
	err := json.NewDecoder(r.Body).Decode(&inputGeometries)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to parse request body: %w", err)
		return
	}

	var geoms []geom.T
	for _, inputGeom := range inputGeometries {
		geometry, err := inputGeom.Decode()
		if err != nil {
			errorHandler <- fmt.Errorf("unable to convert geojson geometry to native type: %w", err)
			return
		}
		geoms = append(geoms, geometry)
	}

	query, err := globals.SqlQueries.Raw("get-withdrawal-rates")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to get query: %w", err)
		return
	}

	var withdrawalRates [][]types.Rate
	for _, geometry := range geoms {
		var rates [][]types.Rate
		err = pgxscan.Select(r.Context(), globals.Db, &rates, query, geometry)
		if err != nil {
			errorHandler <- fmt.Errorf("unable to query withdrawal rates: %w", err)
			return
		}
		withdrawalRates = append(withdrawalRates, rates...)
	}
	if len(withdrawalRates) == 0 {
		errorHandler <- ErrNoWithdrawalRatesAvailable
		return
	}
	var minTakeout, maxTakeout float64
	for _, rates := range withdrawalRates {
		if len(rates) == 1 {
			minTakeout += rates[0].CubicMeterPerYear()
			maxTakeout += rates[0].CubicMeterPerYear()
			continue
		}

		var possibleRates []float64
		for _, rate := range rates {
			possibleRates = append(possibleRates, rate.CubicMeterPerYear())
		}
		minTakeout += slices.Min(possibleRates)
		maxTakeout += slices.Max(possibleRates)
	}
	json.NewEncoder(w).Encode(struct {
		MinimalTakeout float64 `json:"minimalWithdrawal"`
		MaximalTakeout float64 `json:"maximalWithdrawal"`
	}{
		MinimalTakeout: minTakeout,
		MaximalTakeout: maxTakeout,
	})

}
