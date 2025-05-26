package v2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"

	"microservice/internal/db"
	v2 "microservice/types/v2"
)

func UsageLocations(c *gin.Context) {
	var queryParams struct {
		MunicipalityPrefixes []string `form:"in"`
		Active               *bool    `form:"active"`
		Virtual              *bool    `form:"virtual"`
	}
	_ = c.ShouldBindQuery(&queryParams)

	query, err := db.Queries.Raw("get-locations")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var locations []v2.UsageLocation
	err = pgxscan.Select(c, db.Pool(), &locations, query)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	filteredLocations := make([]v2.UsageLocation, 0)
	if queryParams.Active == nil && queryParams.MunicipalityPrefixes == nil && queryParams.Virtual == nil {
		filteredLocations = locations
		goto output
	}

locations:
	for _, location := range locations {
		if queryParams.MunicipalityPrefixes != nil {
			if location.MunicipalArea == nil {
				continue
			}

			for _, requestedKey := range queryParams.MunicipalityPrefixes {
				cleanedKey := fmt.Sprintf("0%d", *location.MunicipalArea.Key)
				if !strings.HasPrefix(cleanedKey, requestedKey) {
					continue locations
				}
			}
		}

		if queryParams.Active != nil {
			if location.Active != queryParams.Active {
				continue
			}
		}

		if queryParams.Virtual != nil && location.Virtual != nil {
			if *location.Virtual != !*queryParams.Virtual {
				continue
			}
		}

		filteredLocations = append(filteredLocations, location)
	}

output:
	featureCollection := geojson.FeatureCollection{
		Features: make([]*geojson.Feature, 0),
		BBox:     geom.NewBounds(geom.XY),
	}

	for _, loc := range filteredLocations {
		feature, _ := loc.ToFeature()
		featureCollection.Features = append(featureCollection.Features, feature)
		featureCollection.BBox.Extend(loc.Geometry)
	}

	encoded, _ := featureCollection.MarshalJSON()

	c.JSON(http.StatusAccepted, json.RawMessage(encoded))
}
