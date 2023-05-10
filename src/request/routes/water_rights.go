package routes

import (
	"encoding/json"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	requestErrors "microservice/request/error"
	"microservice/structs"
	"microservice/vars/globals"
	"microservice/vars/globals/connections"
	"net/http"
)

func WaterRightDetails(w http.ResponseWriter, r *http.Request) {
	l.Debug().Msg("getting water right number from path")
	waterRightNumber := chi.URLParam(r, "waterRight")
	l.Debug().Str("waterRight", waterRightNumber).Msg("got water right number from path")
	l.Info().Str("waterRight", waterRightNumber).Msg("new request for details about water right")

	var waterRightDetails structs.RawWaterRight

	// now query the water right details
	waterRightRow, err := globals.Queries.Query(
		connections.DbConnection,
		"get-water-right-details",
		waterRightNumber)

	if err != nil {
		l.Error().Str("waterRight", waterRightNumber).Err(err).Msg("error while querying information about water right")
		e, _ := requestErrors.WrapInternalError(err)
		requestErrors.SendError(e, w)
		return
	}

	err = scan.Row(&waterRightDetails, waterRightRow)

	var usageLocations []structs.DetailedUsageLocation
	// now query the water right details
	usageLocationRows, err := globals.Queries.Query(
		connections.DbConnection,
		"get-detailed-locations",
		waterRightNumber)

	if err != nil {
		l.Error().Str("waterRight", waterRightNumber).Err(err).Msg("error while querying information about water right")
		e, _ := requestErrors.WrapInternalError(err)
		requestErrors.SendError(e, w)
		return
	}

	err = scan.Rows(&usageLocations, usageLocationRows)

	if err != nil {
		l.Error().Str("waterRight", waterRightNumber).Err(err).Msg("error while parsing information about water right")
		e, _ := requestErrors.WrapInternalError(err)
		requestErrors.SendError(e, w)
		return
	}

	response := waterRightDetails.ToDetailedWaterRight()
	response.Locations = &usageLocations

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		l.Error().Err(err).Msg("unable to send response")
		e, _ := requestErrors.WrapInternalError(err)
		requestErrors.SendError(e, w)
		return
	}
}
