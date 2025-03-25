package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"microservice/internal/db"
	"microservice/types"
)

const (
	locationFilter string = "ST_CONTAINS(ST_COLLECT(ARRAY ((SELECT geom FROM geodata.shapes WHERE key = ANY ($%d)))),ST_TRANSFORM(location, 4326))" //nolint:lll
	stateFilter    string = "active = $1 OR active IS NULL"
	realityFilter  string = `"real" = $1 OR "real" IS NULL`
	queryWhere     string = `%s WHERE (%s)`
	queryAnd       string = `%s AND (%s)`
)

func Locations(c *gin.Context) {
	var filters struct {
		ARS        []string `form:"in"`
		OnlyActive *bool    `form:"is_active"`
		OnlyReal   *bool    `form:"is_real"`
	}

	err := c.ShouldBindQuery(&filters)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var queryArguments []any //nolint:prealloc
	locationFilters := make([]string, len(filters.ARS))
	for idx, key := range filters.ARS {
		locationFilters[idx] = fmt.Sprintf(locationFilter, idx+1)
		queryArguments = append(queryArguments, key)
	}

	query, err := db.Queries.Raw("get-locations")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	if len(filters.ARS) != 0 {
		query = fmt.Sprintf(queryWhere, query, strings.Join(locationFilters, ` OR `))
	}

	if filters.OnlyActive != nil {
		if strings.Contains(query, `WHERE`) {
			query = fmt.Sprintf(queryAnd, query, stateFilter)
		} else {
			query = fmt.Sprintf(queryWhere, query, stateFilter)
		}
		queryArguments = append(queryArguments, *filters.OnlyActive)
	}

	if filters.OnlyReal != nil {
		if strings.Contains(query, `WHERE`) {
			query = fmt.Sprintf(queryAnd, query, realityFilter)
		} else {
			query = fmt.Sprintf(queryWhere, query, realityFilter)
		}
		queryArguments = append(queryArguments, *filters.OnlyReal)
	}

	var locations []types.UsageLocation
	err = scanner.Select(c, db.Pool(), &locations, query, queryArguments...)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	if len(locations) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, locations)

}
