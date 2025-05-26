//go:build release

package router

import (
	"github.com/gin-gonic/gin"
)

// GenerateRouter returns a new [*gin.Engine] which has been configured
// to run in release scenarios.
// This enables security hardening and decreases the default logging level.
func GenerateRouter() (*gin.Engine, error) {
	r := prepareRouter()
	gin.SetMode(gin.ReleaseMode)
	return r, nil

}
