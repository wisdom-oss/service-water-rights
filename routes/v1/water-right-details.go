package v1

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	common "github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	"microservice/types"
)

var (
	errEmptyWaterRightID = common.ServiceError{
		Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
		Status: http.StatusBadRequest,
		Title:  "Empty Water Right ID",
		Detail: "The required water right id has not been transmitted",
	}

	errUnknownWaterRight = common.ServiceError{
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

	var waterRight types.WaterRight
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

	var locations []types.UsageLocation
	err = pgxscan.Select(c, db.Pool(), &locations, query, waterRight.ID)
	if err != nil {
		c.Abort()

		if pgxscan.NotFound(err) {
			goto output
		}

		_ = c.Error(err)
		return
	}

output:

	multipartWriter := multipart.NewWriter(c.Writer)
	c.Header("Content-Type", multipartWriter.FormDataContentType())

	waterRightPart := make(textproto.MIMEHeader)
	waterRightPart.Set("Content-Disposition", `form-data; name="water-right"`)
	waterRightPart.Set("Content-Type", "application/json")
	w, err := multipartWriter.CreatePart(waterRightPart)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_ = json.NewEncoder(w).Encode(waterRight)

	usageLocationsPart := make(textproto.MIMEHeader)
	usageLocationsPart.Set("Content-Disposition", `form-data; name="usage-locations"`)
	usageLocationsPart.Set("Content-Type", "application/json")
	w, err = multipartWriter.CreatePart(usageLocationsPart)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_ = json.NewEncoder(w).Encode(locations)
	multipartWriter.Close()
}
