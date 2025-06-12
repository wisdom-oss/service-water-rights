package router

import (
	"github.com/gin-gonic/gin"

	internal "microservice/internal/router"
	v1Routes "microservice/routes/v1"
	v2Routes "microservice/routes/v2"
)

// Configure generates a new router and adds routes to the router
//
// The router can also be imported during tests, as long as the tests are in a
// separate package.
// If the tests are in the same package (e.g. routes defined in `v3` and tests
// also defined in `v3`) an import cycle exists.
func Configure() (*gin.Engine, error) {
	r, err := internal.GenerateRouter()
	if err != nil {
		return nil, err
	}

	v1 := r.Group("/v1")
	{
		v1.GET("/", v1Routes.UsageLocations)
		v1.GET("/details/:id", v1Routes.WaterRightDetails)

	}

	v2 := r.Group("/v2")
	{
		v2.GET("/", v2Routes.UsageLocations)
		v2.GET("/water-right-details/:id", v2Routes.WaterRightDetails)
	}

	return r, nil
}
