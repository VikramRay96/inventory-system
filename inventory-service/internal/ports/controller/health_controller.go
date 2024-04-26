package controller

import (
	"github.com/gin-gonic/gin"
	"inventory-system/common/pkg/constants"
	"net/http"
	"os"
	"sync"
)

var (
	healthController     *HealthController
	healthControllerOnce sync.Once
)

// health controller with status method
type HealthController struct{}

func NewHealthController() *HealthController {
	if healthController == nil {
		healthControllerOnce.Do(func() {
			healthController = &HealthController{}
		})
	}
	return healthController
}

func (h HealthController) Status() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := os.Getenv(constants.ENV_IS_VERSION)
		c.JSON(http.StatusOK,
			gin.H{"Version": version})
	}

}
