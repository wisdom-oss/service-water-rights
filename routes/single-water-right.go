package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/chi/v5"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	"github.com/wisdom-oss/service-water-rights/globals"
	"github.com/wisdom-oss/service-water-rights/types"
)

// SingleWaterRight returns the water right referenced by the id supplied in the
// request url.
// Only the latest water right will be returned, accompanied by the internal
// ids of the old versions of the specified water right
func SingleWaterRight(w http.ResponseWriter, r *http.Request) {
	// get the error handler and the status channel
	errorHandler := r.Context().Value(errorMiddleware.ChannelName).(chan<- interface{})

	// get the water right number from the url
	waterRightNumber := strings.TrimSpace(chi.URLParam(r, "water-right-nlwkn-id"))

	// validate the water right number
	if waterRightNumber == "" {
		errorHandler <- ErrEmptyWaterRightID
		return
	}

	// now look up the current water right from the association table
	query, err := globals.SqlQueries.Raw("get-current-wr")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve query for water right lookup: %w", err)
		return
	}

	var internalWaterRightID string
	err = pgxscan.Get(r.Context(), globals.Db, &internalWaterRightID, query, waterRightNumber)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve internal id for selected water right: %w", err)
		return
	}

	internalWaterRightID = strings.TrimSpace(internalWaterRightID)
	if internalWaterRightID == "" {
		errorHandler <- ErrNoWaterRightAvailable
		return
	}

	query, err = globals.SqlQueries.Raw("get-water-right")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve query for water right: %w", err)
		return
	}

	var waterRight types.WaterRight
	err = pgxscan.Get(r.Context(), globals.Db, &waterRight, query, internalWaterRightID)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve water right: %w", err)
		return
	}

	query, err = globals.SqlQueries.Raw("get-water-right-usage-locations")
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve query for water right usage locations: %w", err)
		return
	}

	api, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		errorHandler <- fmt.Errorf("unable to prepare query parser: %w", err)
		return
	}
	scanner, err := pgxscan.NewAPI(api)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to create query parser: %w", err)
	}

	var usageLocations []types.UsageLocation
	err = scanner.Select(r.Context(), globals.Db, &usageLocations, query, internalWaterRightID)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to retrieve water right usage locations: %w", err)
		return
	}

	encodedWaterRight, err := json.Marshal(waterRight)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to encode water right: %w", err)
		return
	}

	encodedUsageLocations, err := json.Marshal(usageLocations)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to encode water right: %w", err)
		return
	}

	var outputBuf bytes.Buffer

	multipartWriter := multipart.NewWriter(&outputBuf)
	w.Header().Set("Content-Type", multipartWriter.FormDataContentType())
	waterRightPartHeader := make(textproto.MIMEHeader)
	waterRightPartHeader.Set("Content-Disposition", `form-data; name="water-right"`)
	waterRightPartHeader.Set("Content-Type", "application/json")
	waterRightPart, err := multipartWriter.CreatePart(waterRightPartHeader)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to create multipart field for water right: %w", err)
		return
	}
	_, err = waterRightPart.Write(encodedWaterRight)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to write data to multipart field for water right: %w", err)
		return
	}

	usageLocationsPartHeader := make(textproto.MIMEHeader)
	usageLocationsPartHeader.Set("Content-Disposition", `form-data; name="usage-locations"`)
	usageLocationsPartHeader.Set("Content-Type", "application/json")
	usageLocationsPart, err := multipartWriter.CreatePart(usageLocationsPartHeader)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to create multipart field for usage locations: %w", err)

		return
	}
	_, err = usageLocationsPart.Write(encodedUsageLocations)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to write data to multipart field for usage locations: %w", err)
		return
	}
	multipartWriter.Close()
	_, _ = io.Copy(w, &outputBuf)

}
