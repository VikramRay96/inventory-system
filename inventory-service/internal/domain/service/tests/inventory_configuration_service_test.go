package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	mockClient "inventory-system/inventory-service/internal/adapters/client/mocks"
	mockRepo "inventory-system/inventory-service/internal/adapters/repository/mocks"
	inventoryServiceDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/dto/response_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	serviceImpl "inventory-system/inventory-service/internal/domain/service/impl"
	"testing"
	"time"
)

func init() {
	logger.GetLogger()
}

var mockInventoryConfigurationRepo *mockRepo.MockIInventoryConfigurationRepository
var mockMongoStorageManagerClient *mockClient.MockIMongoStorageManager
var createdTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

var configRequestDto = request_dto.CreateNewConfigurationRequestBody{
	InventoryName:        "Course",
	CreatedBy:            "admin",
	JsonSchema:           map[string]interface{}{},
	InventoryIdentifiers: []string{"name", "course_id"},
}

func TestCreateNewConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryConfigurationRepo = mockRepo.NewMockIInventoryConfigurationRepository(mockController)
	mockMongoStorageManagerClient = mockClient.NewMockIMongoStorageManager(mockController)

	sut := serviceImpl.NewInventoryConfigurationService(mockInventoryConfigurationRepo, mockMongoStorageManagerClient)

	t.Run("TestCreateNewConfiguration_ShouldReturnNilError_WhenNoErrorOccursInRepoAndClient", func(t *testing.T) {
		mockInventoryConfigurationRepo.EXPECT().CreateNewConfiguration(gomock.Any(), configRequestDto).Return(nil)
		mockMongoStorageManagerClient.EXPECT().CreateCollection(gomock.Any(), configRequestDto.InventoryName, configRequestDto.JsonSchema, configRequestDto.InventoryIdentifiers).Return(nil)

		err := sut.CreateNewConfiguration(context.Background(), configRequestDto)
		assert.Nil(t, err)
	})
	t.Run("TestCreateNewConfiguration_ShouldReturnError_WhenErrorReturnedByRepo", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		mockInventoryConfigurationRepo.EXPECT().CreateNewConfiguration(gomock.Any(), configRequestDto).Return(&errDto)
		err := sut.CreateNewConfiguration(context.Background(), configRequestDto)
		assert.Equal(t, errDto.StatusCode, err.StatusCode)
		assert.Equal(t, errDto.Message, err.Message)
	})
	t.Run("TestCreateNewConfiguration_ShouldReturnErrorAndSaveDataInRepo_ErrorReturnedByClient", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		mockInventoryConfigurationRepo.EXPECT().CreateNewConfiguration(gomock.Any(), configRequestDto).Return(nil)
		mockMongoStorageManagerClient.EXPECT().CreateCollection(gomock.Any(), configRequestDto.InventoryName, configRequestDto.JsonSchema, configRequestDto.InventoryIdentifiers).Return(&errDto)
		err := sut.CreateNewConfiguration(context.Background(), configRequestDto)
		assert.Equal(t, errDto.StatusCode, err.StatusCode)
		assert.Equal(t, errDto.Message, err.Message)
	})

}

func TestGetInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryConfigurationRepo = mockRepo.NewMockIInventoryConfigurationRepository(mockController)
	mockMongoStorageManagerClient = mockClient.NewMockIMongoStorageManager(mockController)

	sut := serviceImpl.NewInventoryConfigurationService(mockInventoryConfigurationRepo, mockMongoStorageManagerClient)
	inventoryName := "Course"

	var configRepoResponse = inventoryServiceDto.InventoryConfiguration{
		InventoryName:        "Course",
		InventoryIdentifiers: []string{"name", "course_id"},
		CreatedBy:            "admin",
		CreatedOn:            createdTime,
		JsonSchema:           map[string]interface{}{},
		IsDeleted:            false,
	}
	var serviceResponse = response_dto.InventoryConfigurationResponseDto{
		InventoryName:        "Course",
		InventoryIdentifiers: []string{"name", "course_id"},
		CreatedBy:            "admin",
		CreatedOn:            createdTime,
		JsonSchema:           map[string]interface{}{},
	}

	t.Run("TestGetInventoryConfiguration_ShouldReturnResponseWithNoError_WhenNoErrorOccursAndConfigurationNotDeleted", func(t *testing.T) {
		configRepoResponse.IsDeleted = false
		mockInventoryConfigurationRepo.EXPECT().FetchInventoryConfigurationByName(gomock.Any(), inventoryName).Return(&configRepoResponse, nil)
		resp, err := sut.GetInventoryConfiguration(context.Background(), inventoryName)
		assert.Nil(t, err)
		assert.Equal(t, resp, &serviceResponse)
	})
	t.Run("TestGetInventoryConfiguration_ShouldReturnError_WhenConfigurationIsDeleted", func(t *testing.T) {

		configRepoResponse.IsDeleted = true
		mockInventoryConfigurationRepo.EXPECT().FetchInventoryConfigurationByName(gomock.Any(), inventoryName).Return(&configRepoResponse, nil)
		resp, err := sut.GetInventoryConfiguration(context.Background(), inventoryName)

		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS204)
		assert.Equal(t, err, &expectedErr)
		assert.Equal(t, resp, &serviceResponse)
	})

	t.Run("TestGetInventoryConfiguration_ShouldReturnError_WhenErrorWhileFetchingConfigFromRepo", func(t *testing.T) {

		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		mockInventoryConfigurationRepo.EXPECT().FetchInventoryConfigurationByName(gomock.Any(), inventoryName).Return(nil, &errorDto)
		resp, err := sut.GetInventoryConfiguration(context.Background(), inventoryName)

		assert.Equal(t, err.StatusCode, errorDto.StatusCode)
		assert.Equal(t, err.Message, errorDto.Message)
		assert.Nil(t, resp)
	})

}

func TestDeleteInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryConfigurationRepo = mockRepo.NewMockIInventoryConfigurationRepository(mockController)
	mockMongoStorageManagerClient = mockClient.NewMockIMongoStorageManager(mockController)

	sut := serviceImpl.NewInventoryConfigurationService(mockInventoryConfigurationRepo, mockMongoStorageManagerClient)
	inventoryName := "Course"

	t.Run("TestDeleteInventoryConfiguration_ShouldReturnNoError_WhenNoErrorOccursInRepo", func(t *testing.T) {
		mockInventoryConfigurationRepo.EXPECT().DeleteInventoryConfigurationByName(gomock.Any(), inventoryName).Return(nil)
		err := sut.DeleteInventoryConfiguration(context.Background(), inventoryName)
		assert.Nil(t, err)
	})
	t.Run("TestDeleteInventoryConfiguration_ShouldReturnError_WhenErrorOccursInRepo", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		mockInventoryConfigurationRepo.EXPECT().DeleteInventoryConfigurationByName(gomock.Any(), inventoryName).Return(&errDto)
		err := sut.DeleteInventoryConfiguration(context.Background(), inventoryName)
		assert.Equal(t, err.StatusCode, errDto.StatusCode)
		assert.Equal(t, err.Message, errDto.Message)
	})

}

func TestGetAllInventoryConfiguration(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryConfigurationRepo = mockRepo.NewMockIInventoryConfigurationRepository(mockController)
	mockMongoStorageManagerClient = mockClient.NewMockIMongoStorageManager(mockController)

	var configRepoResponse1 = inventoryServiceDto.InventoryConfiguration{
		InventoryName:        "Course",
		InventoryIdentifiers: []string{"name", "course_id"},
		CreatedBy:            "admin",
		CreatedOn:            createdTime,
		JsonSchema:           map[string]interface{}{},
		IsDeleted:            false,
	}
	var configRepoResponse2 = inventoryServiceDto.InventoryConfiguration{
		InventoryName:        "Placement",
		InventoryIdentifiers: []string{"name", "placement_id"},
		CreatedBy:            "admin",
		CreatedOn:            createdTime,
		JsonSchema:           map[string]interface{}{},
		IsDeleted:            true,
	}
	var configListRepo []inventoryServiceDto.InventoryConfiguration
	configListRepo = append(configListRepo, configRepoResponse1)
	configListRepo = append(configListRepo, configRepoResponse2)

	sut := serviceImpl.NewInventoryConfigurationService(mockInventoryConfigurationRepo, mockMongoStorageManagerClient)

	t.Run("TestTestGetAllInventoryConfiguration_ShouldReturnNoError_WhenNoErrorOccursInRepo", func(t *testing.T) {
		//Only 1 active config in list of 2

		mockInventoryConfigurationRepo.EXPECT().FetchAllInventoryConfiguration(gomock.Any()).Return(configListRepo, nil)
		resp, err := sut.GetAllInventoryConfiguration(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, len(resp), 1)

	})
	t.Run("TestTestGetAllInventoryConfiguration_ShouldReturnError_WhenErrorOccursInRepo", func(t *testing.T) {
		var errDto dto.ErrorResponseDto
		errDto.SetError(status_code.IMS500)

		mockInventoryConfigurationRepo.EXPECT().FetchAllInventoryConfiguration(gomock.Any()).Return(nil, &errDto)
		resp, err := sut.GetAllInventoryConfiguration(context.Background())
		assert.Equal(t, err.StatusCode, errDto.StatusCode)
		assert.Equal(t, err.Message, errDto.Message)
		assert.Nil(t, resp)
	})

}
