package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v3/middleware/gin/jwt"

	"microservice/internal"
	internalRouter "microservice/internal/router"
	v1Routes "microservice/routes/v1"
)

// Configure generates a new router and adds routes to the router
//
// The router can also be imported during tests, as long as the tests are in a
// separate package.
// If the tests are in the same package (e.g. routes defined in `v3` and tests
// also defined in `v3`) an import cycle exists.
func Configure() (*gin.Engine, error) {
	r, err := internalRouter.GenerateRouter()
	if err != nil {
		return nil, err
	}
	r.Use(gzip.Gzip(gzip.BestCompression))

	scopeRequirer := jwt.ScopeRequirer{}
	scopeRequirer.Configure(internal.ServiceName)

	v1 := r.Group("", scopeRequirer.RequireRead)
	{
		v1.GET("/", v1Routes.Locations)
		v1.GET("/details/:nlwkn-water-right-id", v1Routes.WaterRightDetails)
		v1.POST("/average-withdrawals", v1Routes.CalculateWaterWithdrawal)
	}

	return r, nil
}
