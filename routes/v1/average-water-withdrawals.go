package v1

import (
	"net/http"
	"slices"
	"sync"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"golang.org/x/sync/errgroup"

	wisdom "github.com/wisdom-oss/common-go/v3/types"

	"microservice/globals"
	"microservice/internal/db"
	"microservice/types"
)

var (
	errInvalidBody = wisdom.ServiceError{
		Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.4.1",
		Status: http.StatusBadRequest,
		Title:  "Invalid Request Body",
		Detail: "Please only transmit GeoJSON geometries as request body",
	}
)

func AverageWaterTakeout(c *gin.Context) {
	var geometries []geojson.Geometry
	if err := c.ShouldBindBodyWithJSON(&geometries); err != nil {
		c.Abort()
		errInvalidBody.Emit(c)
		return
	}

	query, err := globals.SqlQueries.Raw("get-withdrawal-rates")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var withdrawalRates [][]types.Rate
	var lock sync.Mutex
	var paralel errgroup.Group
	for _, geometry := range geometries {
		paralel.Go(func() error {
			decodedGeometry, err := geometry.Decode()
			if err != nil {
				return err
			}

			decodedGeometry, err = geom.SetSRID(decodedGeometry, 4326) //nolint:mnd
			if err != nil {
				return err
			}

			var r [][]types.Rate
			err = pgxscan.Select(c, db.Pool(), &r, query, decodedGeometry)
			if err != nil {
				if pgxscan.NotFound(err) {
					return nil
				}
				return err
			}

			lock.Lock()
			withdrawalRates = append(withdrawalRates, r...)
			lock.Unlock()

			return nil

		})
	}

	err = paralel.Wait()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var minimalTakeout, maximalTakeout float64
	for _, rates := range withdrawalRates {
		if len(rates) == 1 {
			minimalTakeout += rates[0].CubicMeterPerYear()
			maximalTakeout += rates[0].CubicMeterPerYear()
		}

		var possibleRates []float64
		for _, rate := range rates {
			possibleRates = append(possibleRates, rate.CubicMeterPerYear())
		}

		minimalTakeout += slices.Min(possibleRates)
		maximalTakeout += slices.Max(possibleRates)
	}

	var takeout struct {
		Minimal float64 `json:"minimalWithdrawal"`
		Maximal float64 `json:"maximalWithdrawal"`
	}
	takeout.Maximal = maximalTakeout
	takeout.Minimal = minimalTakeout

	c.JSON(http.StatusOK, takeout)

}
