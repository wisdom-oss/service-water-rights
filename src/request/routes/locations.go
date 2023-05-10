package routes

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	geojson "github.com/paulmach/go.geojson"
	requestErrors "microservice/request/error"
	"microservice/structs"
	"microservice/vars/globals"
	"microservice/vars/globals/connections"
	"net/http"
	"strconv"
)

var l = globals.HttpLogger

/*
UsageLocations

This handler shows how a basic handler works and how to send back a message
*/
func UsageLocations(w http.ResponseWriter, request *http.Request) {
	l.Info().Msg("new request for usage locations")
	// create some non-assigned booleans for building the request filtering
	var areaFilterSet, realityFilterSet, stateFilterSet bool
	var areaKeys []string
	var realRightStrings, activeRightStrings []string
	var isRealRight, isActiveRight bool

	areaKeys, areaFilterSet = request.URL.Query()["in"]
	realRightStrings, realityFilterSet = request.URL.Query()["is_real"]
	if realityFilterSet {
		if b, err := strconv.ParseBool(realRightStrings[0]); err == nil {
			isRealRight = b
		} else {
			l.Warn().Msg("reality filter not set with an appropriate boolean value. Disabling filter")
			realityFilterSet = false
		}
	}
	activeRightStrings, stateFilterSet = request.URL.Query()["is_real"]
	if stateFilterSet {
		if b, err := strconv.ParseBool(activeRightStrings[0]); err == nil {
			isActiveRight = b
		} else {
			l.Warn().Msg("state filter not set with an appropriate boolean value. Disabling filter")
			stateFilterSet = false
		}
	}

	var resultRows *sql.Rows
	var queryError error

	switch {
	case areaFilterSet && realityFilterSet && stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"area", "reality", "state"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-reality-state-and-location",
			isRealRight, isActiveRight, pq.Array(areaKeys),
		)
		break
	case !areaFilterSet && realityFilterSet && stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"reality", "state"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-reality-and-state",
			isRealRight, isActiveRight,
		)
		break
	case areaFilterSet && !realityFilterSet && stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"area", "state"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-state-and-location",
			isActiveRight, pq.Array(areaKeys),
		)
		break
	case areaFilterSet && realityFilterSet && !stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"area", "reality"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-reality-and-location",
			isRealRight, pq.Array(areaKeys),
		)
		break
	case !areaFilterSet && !realityFilterSet && stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"state"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-state",
			isActiveRight,
		)
		break
	case !areaFilterSet && realityFilterSet && !stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"reality"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-reality",
			isRealRight,
		)
		break
	case areaFilterSet && !realityFilterSet && !stateFilterSet:
		l.Info().Strs("enabledFilters", []string{"area"}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-water-rights-by-location",
			areaKeys,
		)
		break
	case !areaFilterSet && !realityFilterSet && !stateFilterSet:
		l.Warn().Strs("enabledFilters", []string{}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-all-water-rights",
		)
		break
	default:
		l.Warn().Strs("enabledFilters", []string{}).Msg("querying database for usage locations")
		resultRows, queryError = globals.Queries.Query(
			connections.DbConnection,
			"get-all-water-rights",
		)
		break
	}

	if queryError != nil {
		l.Error().Err(queryError).Msg("database query failed")
		e, _ := requestErrors.WrapInternalError(queryError)
		requestErrors.SendError(e, w)
		return
	}

	defer func(r *sql.Rows) {
		err := r.Close()
		if err != nil {
			l.Warn().Err(err).Msg("error while closing the active result rows")
		}
	}(resultRows)

	var locations []structs.UsageLocation

	for resultRows.Next() {
		var isReal, isActive *bool
		var internalId, externalId int
		var name *string
		var location *geojson.Geometry

		err := resultRows.Scan(&internalId, &externalId, &isActive, &isReal, &name, &location)
		if err != nil {
			l.Error().Err(err).Msg("unable to parse database rows into response object")
			e, _ := requestErrors.WrapInternalError(err)
			requestErrors.SendError(e, w)
			return
		}

		locations = append(locations, structs.UsageLocation{
			ID:         internalId,
			WaterRight: externalId,
			IsActive:   isActive,
			IsReal:     isReal,
			Name:       name,
			Location:   location,
		})
	}

	if len(locations) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(locations)
	if err != nil {
		l.Error().Err(err).Msg("unable to send response")
		e, _ := requestErrors.WrapInternalError(err)
		requestErrors.SendError(e, w)
		return
	}
}
