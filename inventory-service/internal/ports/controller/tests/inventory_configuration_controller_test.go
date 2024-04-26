package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/dto/response_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	mockServices "inventory-system/inventory-service/internal/domain/service/mocks"
	"inventory-system/inventory-service/internal/ports/controller"
	"inventory-system/inventory-service/internal/ports/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func init() {
	logger.GetLogger()
}

var inventoryConfigurationServiceMock *mockServices.MockIInventoryConfigurationService

var getInventoryConfigurationResponseDto = response_dto.InventoryConfigurationResponseDto{
	InventoryName:        "Course",
	InventoryIdentifiers: []string{"name", "course_id"},
	CreatedBy:            "admin",
	CreatedOn:            time.Now(),
	JsonSchema:           map[string]interface{}{},
}

var createNewInventoryConfigurationRequestDto = request_dto.CreateNewConfigurationRequestBody{
	InventoryName:        "Course",
	CreatedBy:            "admin",
	JsonSchema:           map[string]interface{}{},
	InventoryIdentifiers: []string{"name", "course_id"},
}

func SetupInventoryConfigurationRouter(mockController *gomock.Controller) *gin.Engine {
	router := gin.New()
	inventoryConfigurationServiceMock = mockServices.NewMockIInventoryConfigurationService(mockController)
	validator := utils.NewRequestValidator()
	inventoryConfigurationController := controller.NewInventoryConfigurationController(inventoryConfigurationServiceMock, validator)
	InventoryServiceRouter := router.Group("inventory-service/api/v1")
	inventoryConfigurations := InventoryServiceRouter.Group("/inventory/configurations")
	{
		inventoryConfigurations.GET("/:inventoryName", inventoryConfigurationController.GetConfiguration())
		inventoryConfigurations.POST("", inventoryConfigurationController.CreateNewConfiguration())
		inventoryConfigurations.GET("", inventoryConfigurationController.GetAllConfiguration())
		inventoryConfigurations.DELETE("/:inventoryName", inventoryConfigurationController.DeleteConfiguration())
	}
	return router
}

func TestGetInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryConfigurationRouter(mockController)

	inventoryName := "Course"
	url := "/inventory-service/api/v1/inventory/configurations/:inventoryName"
	url = strings.Replace(url, ":inventoryName", inventoryName, -1)

	t.Run("TestGetInventoryConfiguration_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		inventoryConfigurationServiceMock.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&getInventoryConfigurationResponseDto, nil)
		req, _ := http.NewRequest("GET", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestGetInventoryConfiguration_ShouldReturnStatus500_WhenInventoryConfigurationServiceReturnsError500", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		inventoryConfigurationServiceMock.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(nil, &errorDto)
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

func TestCreateNewInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryConfigurationRouter(mockController)

	url := "/inventory-service/api/v1/inventory/configurations"

	t.Run("TestCreateNewInventoryConfiguration_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		inventoryConfigurationServiceMock.EXPECT().CreateNewConfiguration(gomock.Any(), createNewInventoryConfigurationRequestDto).Return(nil)

		body, _ := json.Marshal(createNewInventoryConfigurationRequestDto)
		reqBody := string(body)
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestCreateNewInventoryConfiguration_ShouldReturnsError_WhenInventoryServiceReturnsError", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		inventoryConfigurationServiceMock.EXPECT().CreateNewConfiguration(gomock.Any(), createNewInventoryConfigurationRequestDto).Return(&errDto)
		body, _ := json.Marshal(createNewInventoryConfigurationRequestDto)
		reqBody := string(body)
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).Message, responseValue.Message)
	})
	t.Run("TestCreateNewInventoryConfiguration_ShouldReturnStatus400_WhenValidationFails", func(t *testing.T) {
		//Inventory Identifiers is required cant be nil
		createNewInventoryConfigurationRequestDto.InventoryIdentifiers = nil

		body, _ := json.Marshal(createNewInventoryConfigurationRequestDto)
		reqBody := string(body)
		req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS400).StatusCode, responseValue.StatusCode)
		assert.Equal(t, "inventory_identifiers : Field is missing", responseValue.Message)
	})
	t.Run("TestCreateNewInventoryConfiguration_ShouldReturnStatus400_WhenJsonInvalid", func(t *testing.T) {
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

func TestDeleteInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryConfigurationRouter(mockController)

	inventoryName := "Course"
	url := "/inventory-service/api/v1/inventory/configurations/:inventoryName"
	url = strings.Replace(url, ":inventoryName", inventoryName, -1)

	t.Run("TestDeleteInventoryConfiguration_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		inventoryConfigurationServiceMock.EXPECT().DeleteInventoryConfiguration(gomock.Any(), inventoryName).Return(nil)
		req, _ := http.NewRequest("DELETE", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestDeleteInventoryConfiguration_ShouldReturnStatus500_WhenInventoryConfigurationServiceReturnsError500", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		inventoryConfigurationServiceMock.EXPECT().DeleteInventoryConfiguration(gomock.Any(), inventoryName).Return(&errorDto)
		req, _ := http.NewRequest("DELETE", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).Message, responseValue.Message)
	})
}

func TestGetAllInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()
	router := SetupInventoryConfigurationRouter(mockController)

	inventoryName := "Course"
	url := "/inventory-service/api/v1/inventory/configurations"
	url = strings.Replace(url, ":inventoryName", inventoryName, -1)

	t.Run("TestTestGetAllInventoryConfiguration_ShouldReturnStatus200_WhenNoErrorOccurs", func(t *testing.T) {
		var inventoryConfigurations []response_dto.InventoryConfigurationResponseDto

		inventoryConfigurationServiceMock.EXPECT().GetAllInventoryConfiguration(gomock.Any()).Return(inventoryConfigurations, nil)
		req, _ := http.NewRequest("GET", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS200).Message, responseValue.Message)
	})
	t.Run("TestDeleteInventoryConfiguration_ShouldReturnStatus500_WhenInventoryConfigurationServiceReturnsError500", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		inventoryConfigurationServiceMock.EXPECT().GetAllInventoryConfiguration(gomock.Any()).Return(nil, &errorDto)
		req, _ := http.NewRequest("GET", url, nil)
		recordedResponse := httptest.NewRecorder()
		router.ServeHTTP(recordedResponse, req)
		var responseValue dto.ResponseDto
		_ = json.Unmarshal(recordedResponse.Body.Bytes(), &responseValue)

		assert.Equal(t, http.StatusOK, recordedResponse.Code)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).StatusCode, responseValue.StatusCode)
		assert.Equal(t, dto.GetStatusDetails(status_code.IMS500).Message, responseValue.Message)
		assert.Equal(t, nil, responseValue.Data)
	})
}
