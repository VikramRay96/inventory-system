package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"inventory-system/common/pkg/logger"
	portConstants "inventory-system/inventory-service/internal/ports/constants"
	"inventory-system/inventory-service/internal/ports/controller"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	logger.GetLogger()
}

func healthSetupRouter() *gin.Engine {
	router := gin.New()
	healthController := controller.NewHealthController()
	router.GET(portConstants.HEALTH_CHECK_PATH, healthController.Status())
	return router

}
func TestHealthCheck(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := healthSetupRouter()

	t.Run("TestGetInventory_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		url := "/inventory-service/health"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

}
