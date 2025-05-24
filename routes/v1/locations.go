package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	"microservice/internal/db"
	"microservice/types"
)

func UsageLocations(c *gin.Context) {
	var queryParams struct {
		MunicipalityKeys []string `form:"in"`
		Active           *bool    `form:"is_active"`
		Real             *bool    `form:"is_real"`
	}
	_ = c.ShouldBindQuery(&queryParams)

	query, err := db.Queries.Raw("get-locations")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var usageLocations []types.UsageLocation
	err = pgxscan.Select(c, db.Pool(), &usageLocations, query)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	filteredLocations := make([]types.UsageLocation, 0)

	if queryParams.Active == nil && queryParams.MunicipalityKeys == nil && queryParams.Real == nil {
		filteredLocations = usageLocations
		goto output
	}

locations:
	for _, location := range usageLocations {
		if queryParams.MunicipalityKeys != nil {
			if location.MunicipalArea == nil {
				continue
			}

			for _, requestedKey := range queryParams.MunicipalityKeys {
				cleanedKey := fmt.Sprintf("0%d", location.MunicipalArea.Key.Int)
				fmt.Println(cleanedKey, requestedKey, len(cleanedKey)-len(requestedKey))
				if !strings.HasPrefix(cleanedKey, requestedKey) {
					continue locations
				}
				fmt.Println("match found", cleanedKey, requestedKey, len(cleanedKey)-len(requestedKey))
			}
		}

		if queryParams.Active != nil {
			if location.Active.Bool != *queryParams.Active {
				continue
			}
		}

		if queryParams.Real != nil {
			if location.Real.Bool != *queryParams.Real {
				continue
			}
		}

		filteredLocations = append(filteredLocations, location)
	}
output:
	c.JSON(http.StatusOK, filteredLocations)

}
