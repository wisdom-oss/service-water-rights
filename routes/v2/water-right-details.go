package v2

import (
	"net/http"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	v2 "microservice/types/v2"
)

var (
	errEmptyWaterRightID = types.ServiceError{
		Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
		Status: http.StatusBadRequest,
		Title:  "Empty Water Right ID",
		Detail: "The required water right id has not been transmitted",
	}

	errUnknownWaterRight = types.ServiceError{
		Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
		Status: http.StatusNotFound,
		Title:  "Unknown Water Right",
		Detail: "The specified water right is not stored in the database",
	}
)

func WaterRightDetails(c *gin.Context) {
	waterRightID := strings.TrimSpace(c.Param("id"))
	if waterRightID == "" {
		c.Abort()
		errEmptyWaterRightID.Emit(c)
		return
	}

	query, err := db.Queries.Raw("v2_get-water-right")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var waterRight v2.WaterRight
	err = pgxscan.Get(c, db.Pool(), &waterRight, query, waterRightID)
	if err != nil {
		c.Abort()

		if pgxscan.NotFound(err) {
			errUnknownWaterRight.Emit(c)
			return
		}

		_ = c.Error(err)
		return
	}

	query, err = db.Queries.Raw("v2_get-water-right-usage-locations")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var locations []v2.UsageLocation
	err = pgxscan.Select(c, db.Pool(), &locations, query, waterRight.Identifiers.Database)
	if err != nil {
		c.Abort()

		if pgxscan.NotFound(err) {
			waterRight.AssociatedUsageLocations = nil
			goto output
		}

		_ = c.Error(err)
		return
	}

output:
	waterRight.AssociatedUsageLocations = locations
	c.JSON(http.StatusOK, waterRight)
}
