package v1

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	commonTypes "github.com/wisdom-oss/common-go/v3/types"

	"microservice/internal/db"
	"microservice/types"
)

var ErrEmptyWaterRightID = commonTypes.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Empty Water Right ID",
	Detail: "The water right id set in the query is empty. Please check your query",
}

var ErrUnknownWaterRight = commonTypes.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: http.StatusNotFound,
	Title:  "Unkown Water Right ID",
	Detail: "The supplied water right id does not exist",
}

func WaterRightDetails(c *gin.Context) {
	waterRightNumber := strings.TrimSpace(c.Param("nlwkn-water-right-id"))
	if waterRightNumber == "" {
		c.Abort()
		ErrEmptyWaterRightID.Emit(c)
		return
	}

	query, err := db.Queries.Raw("get-current-wr")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var internalID string
	err = pgxscan.Get(c, db.Pool(), &internalID, query, waterRightNumber)
	if err != nil {
		c.Abort()
		if pgxscan.NotFound(err) {
			ErrUnknownWaterRight.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	query, err = db.Queries.Raw("get-water-right")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var waterRight types.WaterRight
	err = scanner.Get(c, db.Pool(), &waterRight, query, internalID)
	if err != nil {
		c.Abort()
		if pgxscan.NotFound(err) {
			ErrUnknownWaterRight.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	query, err = db.Queries.Raw("get-water-right-usage-locations")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var usageLocations []types.UsageLocation
	err = scanner.Select(c, db.Pool(), &usageLocations, query, internalID)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)

	mpWriter := multipart.NewWriter(c.Writer)
	c.Header("Content-Type", mpWriter.FormDataContentType())

	wrPartHeader := make(textproto.MIMEHeader)
	wrPartHeader.Set("Content-Disposition", `form-data; name="water-right"`)
	wrPartHeader.Set("Content-Type", "application/json")
	wrPart, err := mpWriter.CreatePart(wrPartHeader)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	err = json.NewEncoder(wrPart).Encode(waterRight)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	locPartHeader := make(textproto.MIMEHeader)
	locPartHeader.Set("Content-Disposition", `form-data; name="usage-locaions"`)
	locPartHeader.Set("Content-Type", "application/json")
	locPart, err := mpWriter.CreatePart(locPartHeader)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	err = json.NewEncoder(locPart).Encode(&usageLocations)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_ = mpWriter.Close()

}
