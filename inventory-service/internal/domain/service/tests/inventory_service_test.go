package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	mockRepo "inventory-system/inventory-service/internal/adapters/repository/mocks"
	"inventory-system/inventory-service/internal/common/dto/response_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	serviceImpl "inventory-system/inventory-service/internal/domain/service/impl"
	mockServices "inventory-system/inventory-service/internal/domain/service/mocks"
	"testing"
)

func init() {
	logger.GetLogger()
}

var mockInventoryRepo *mockRepo.MockIInventoryRepository
var mockInventoryConfigurationService *mockServices.MockIInventoryConfigurationService

func TestCreateNewInventory(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryRepo = mockRepo.NewMockIInventoryRepository(mockController)
	mockInventoryConfigurationService = mockServices.NewMockIInventoryConfigurationService(mockController)

	sut := serviceImpl.NewInventoryService(mockInventoryRepo, mockInventoryConfigurationService)
	inventoryName := "Course"
	var item interface{}
	var serviceResponse = response_dto.InventoryConfigurationResponseDto{
		InventoryName:        "Course",
		InventoryIdentifiers: []string{"name", "course_id"},
		CreatedBy:            "admin",
		CreatedOn:            createdTime,
		JsonSchema:           map[string]interface{}{},
	}

	t.Run("TestCreateNewInventory_ShouldReturnNilError_WhenNoErrorOccursInFetchConfigAndCreateItem", func(t *testing.T) {
		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, nil)
		mockInventoryRepo.EXPECT().CreateNewInventoryGivenInventoryName(gomock.Any(), item, inventoryName).Return(nil)
		err := sut.CreateNewInventory(context.Background(), item, inventoryName)

		assert.Nil(t, err)
	})
	t.Run("TestCreateNewInventory_ShouldReturnError_WhenNoErrorOccursInFetchConfig", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS204)

		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(nil, &errorDto)
		err := sut.CreateNewInventory(context.Background(), item, inventoryName)

		assert.Equal(t, err.StatusCode, errorDto.StatusCode)
		assert.Equal(t, err.Message, errorDto.Message)
	})
	t.Run("TestCreateNewInventory_ShouldReturnError_WhenNoErrorOccursInCreateItem", func(t *testing.T) {
		var errorDto dto.ErrorResponseDto
		errorDto.SetError(status_code.IMS500)

		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, nil)
		mockInventoryRepo.EXPECT().CreateNewInventoryGivenInventoryName(gomock.Any(), item, inventoryName).Return(&errorDto)
		err := sut.CreateNewInventory(context.Background(), item, inventoryName)

		assert.Equal(t, err.StatusCode, errorDto.StatusCode)
		assert.Equal(t, err.Message, errorDto.Message)
	})

}
func TestGetInventory(t *testing.T) {
	var mockController = gomock.NewController(t)
	defer mockController.Finish()

	mockInventoryRepo = mockRepo.NewMockIInventoryRepository(mockController)
	mockInventoryConfigurationService = mockServices.NewMockIInventoryConfigurationService(mockController)

	sut := serviceImpl.NewInventoryService(mockInventoryRepo, mockInventoryConfigurationService)
	inventoryName := "Course"

	item := bson.M{"course_name": "DSA", "course_id": "123456", "duration": 2}

	t.Run("TestGetInventory_ShouldReturnNilError_WhenNoErrorOccurs", func(t *testing.T) {
		var serviceResponse = response_dto.InventoryConfigurationResponseDto{
			InventoryName:        "Course",
			InventoryIdentifiers: []string{"course_name", "course_id"},
			CreatedBy:            "admin",
			CreatedOn:            createdTime,
			JsonSchema:           map[string]interface{}{},
		}
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")
		item = bson.M{"course_name": "DSA", "course_id": "123456", "duration": 2}

		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, nil)
		mockInventoryRepo.EXPECT().FetchInventory(gomock.Any(), inventoryName, gomock.Any(), gomock.Any()).Return(item, nil)
		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)

		assert.Nil(t, err)
		assert.Equal(t, resp, item)
	})
	t.Run("TestGetInventory_ShouldReturnError_WhenMoreThan1AttributeInFilterMap", func(t *testing.T) {
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")

		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS112)
		filterMap["course_id"] = append(filterMap["course_id"], "123")
		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)
		assert.Nil(t, resp)
		assert.Equal(t, err.StatusCode, expectedErr.StatusCode)
		assert.Equal(t, err.Message, expectedErr.Message)
	})
	t.Run("TestGetInventory_ShouldReturnError_WhenConfigurationIsDeleted", func(t *testing.T) {
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")
		item = bson.M{"course_name": "DSA", "course_id": "123456", "duration": 2}

		var serviceResponse = response_dto.InventoryConfigurationResponseDto{
			InventoryName:        "Course",
			InventoryIdentifiers: []string{"course_name", "course_id"},
			CreatedBy:            "admin",
			CreatedOn:            createdTime,
			JsonSchema:           map[string]interface{}{},
		}
		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS204)
		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, &expectedErr)

		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)
		assert.Nil(t, resp)
		assert.Equal(t, err.StatusCode, expectedErr.StatusCode)
		assert.Equal(t, err.Message, expectedErr.Message)
	})
	t.Run("TestGetInventory_ShouldReturnError_WhenConfigurationReturnError", func(t *testing.T) {
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")
		item = bson.M{"course_name": "DSA", "course_id": "123456", "duration": 2}

		var serviceResponse = response_dto.InventoryConfigurationResponseDto{
			InventoryName:        "Course",
			InventoryIdentifiers: []string{"course_name", "course_id"},
			CreatedBy:            "admin",
			CreatedOn:            createdTime,
			JsonSchema:           map[string]interface{}{},
		}
		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS500)
		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, &expectedErr)

		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)
		assert.Nil(t, resp)
		assert.Equal(t, err.StatusCode, expectedErr.StatusCode)
		assert.Equal(t, err.Message, expectedErr.Message)
	})

	t.Run("TestGetInventory_ShouldReturnError_WhenFilterAttributeNotInventoryIdentifierList", func(t *testing.T) {
		var serviceResponse = response_dto.InventoryConfigurationResponseDto{
			InventoryName:        "Course",
			InventoryIdentifiers: []string{"course_name", "course_id"},
			CreatedBy:            "admin",
			CreatedOn:            createdTime,
			JsonSchema:           map[string]interface{}{},
		}
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")

		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS113)

		serviceResponse.InventoryIdentifiers = []string{"dummy"}

		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, nil)
		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)
		assert.Nil(t, resp)
		assert.Equal(t, err.StatusCode, expectedErr.StatusCode)
		assert.Equal(t, err.Message, expectedErr.Message)
	})
	t.Run("TestGetInventory_ShouldReturnError_WhenFetchFromRepoReturnError", func(t *testing.T) {
		var serviceResponse = response_dto.InventoryConfigurationResponseDto{
			InventoryName:        "Course",
			InventoryIdentifiers: []string{"course_name", "course_id"},
			CreatedBy:            "admin",
			CreatedOn:            createdTime,
			JsonSchema:           map[string]interface{}{},
		}
		filterMap := make(map[string][]string)
		filterMap["course_name"] = append(filterMap["course_name"], "123")
		item = bson.M{"course_name": "DSA", "course_id": "123456", "duration": 2}

		var expectedErr dto.ErrorResponseDto
		expectedErr.SetError(status_code.IMS500)

		mockInventoryConfigurationService.EXPECT().GetInventoryConfiguration(gomock.Any(), inventoryName).Return(&serviceResponse, nil)
		mockInventoryRepo.EXPECT().FetchInventory(gomock.Any(), inventoryName, gomock.Any(), gomock.Any()).Return(item, &expectedErr)

		resp, err := sut.GetInventory(context.Background(), inventoryName, filterMap)
		assert.Nil(t, resp)
		assert.Equal(t, err.StatusCode, expectedErr.StatusCode)
		assert.Equal(t, err.Message, expectedErr.Message)
	})

}
