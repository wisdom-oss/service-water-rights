package v1

import (
	"encoding/json"
	"net/http"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geom/encoding/geojson"

	"microservice/internal/db"
	"microservice/types"
)

func CalculateWaterWithdrawal(c *gin.Context) {
	var geometries []geojson.Geometry
	err := c.ShouldBindBodyWithJSON(&geometries)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	query, err := db.Queries.Raw("get-withdrawal-rates")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var rates [][]types.Rate
	var errors []error
	var wg sync.WaitGroup
	var errormutex sync.Mutex
	var ratemutex sync.Mutex

	for _, geometry := range geometries {
		wg.Add(1)
		go func(geom geojson.Geometry) {
			defer wg.Done()
			var shapeRates [][]types.Rate

			param, _ := json.Marshal(geom)
			err := scanner.Select(c, db.Pool(), &shapeRates, query, param)
			if err != nil {
				errormutex.Lock()
				errors = append(errors, err)
				errormutex.Unlock()
			}

			ratemutex.Lock()
			rates = append(rates, rates...)
			ratemutex.Unlock()
		}(geometry)
	}

	wg.Wait()

	var takeout struct {
		Min float64 `json:"minimalWithdrawal"`
		Max float64 `json:"maximalWithdrawal"`
	}
	for _, rates := range rates {
		if len(rates) == 1 {
			takeout.Min += rates[0].CubicMeterPerYear()
			takeout.Max += rates[0].CubicMeterPerYear()
			continue
		}

		var singleRates []float64
		for _, rate := range rates {
			singleRates = append(singleRates, rate.CubicMeterPerYear())
		}

		takeout.Min += slices.Min(singleRates)
		takeout.Max += slices.Max(singleRates)
	}

	c.JSON(http.StatusOK, takeout)
}
