package pkg

import (
	"github.com/gin-gonic/gin"
	"inventory-system/inventory-service/internal/ports/router"
)

func RouterDriver() *gin.Engine {
	return router.NewRouter()
}
