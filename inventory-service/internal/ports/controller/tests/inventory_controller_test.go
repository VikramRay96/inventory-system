package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/inventory-service/internal/common/status_code"
	mockServices "inventory-system/inventory-service/internal/domain/service/mocks"
	"inventory-system/inventory-service/internal/ports/controller"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	logger.GetLogger()
}

var inventoryServiceMock *mockServices.MockIInventoryService

func SetupInventoryRouter(mockController *gomock.Controller) *gin.Engine {
	router := gin.New()
	inventoryServiceMock = mockServices.NewMockIInventoryService(mockController)
	inventoryController := controller.NewInventoryController(inventoryServiceMock)
	InventoryServiceRouter := router.Group("inventory-service/api/v1")
	inventory := InventoryServiceRouter.Group("/inventory")
	{
		inventory.POST("/:inventoryName", inventoryController.AddNewInventory())
		inventory.GET("/:inventoryName", inventoryController.GetInventory())
	}
	return router
}

func TestGetInventory(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryRouter(mockController)

	inventoryName := "Course"
	url := "/inventory-service/api/v1/inventory/:inventoryName?name=123"
	url = strings.Replace(url, ":inventoryName", inventoryName, -1)

	t.Run("TestGetInventory_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		var item bson.M

		inventoryServiceMock.EXPECT().GetInventory(gomock.Any(), inventoryName, gomock.Any()).Return(item, nil)
		req, _ := http.NewRequest("GET", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestGetInventory_ShouldReturnStatus500_WhenInventoryConfigurationServiceReturnsError500", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		inventoryServiceMock.EXPECT().GetInventory(gomock.Any(), inventoryName, gomock.Any()).Return(nil, &errorDto)
		req, _ := http.NewRequest("GET", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).Message, responseValue.Message)
	})
}

func TestCreateNewInventory(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryRouter(mockController)

	inventoryName := "Course"
	url := "/inventory-service/api/v1/inventory/:inventoryName?name=123"
	url = strings.Replace(url, ":inventoryName", inventoryName, -1)

	t.Run("TestCreateNewInventory_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		inventoryServiceMock.EXPECT().CreateNewInventory(gomock.Any(), gomock.Any(), inventoryName).Return(nil)
		reqBody := `{"course_name":"DSA", "course_id" : "123"}`
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestCreateNewInventory_ShouldReturnsError_WhenInventoryServiceReturnsError", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		inventoryServiceMock.EXPECT().CreateNewInventory(gomock.Any(), gomock.Any(), inventoryName).Return(&errDto)
		reqBody := `{"course_name":"DSA", "course_id" : "123"}`
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).Message, responseValue.Message)
	})

	t.Run("TestCreateNewInventory_ShouldReturnStatus400_WhenJsonInvalid", func(t *testing.T) {
		//Request body is nil
		req, _ := http.NewRequest("POST", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS400).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS400).Message, responseValue.Message)
	})

}
