package impl

import (
	"context"
	"inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/common/pkg/utils"
	"inventory-system/inventory-service/internal/adapters/client"
	"inventory-system/inventory-service/internal/adapters/repository"
	inventoryServiceDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/dto/response_dto"
	"inventory-system/inventory-service/internal/common/status_code"
)

type InventoryConfigurationService struct {
	InventoryConfigurationRepository repository.IInventoryConfigurationRepository
	MongoStorageManagerClient        client.IMongoStorageManager
}

func NewInventoryConfigurationService(inventoryConfigurationRepository repository.IInventoryConfigurationRepository, mongoStorageManager client.IMongoStorageManager) *InventoryConfigurationService {
	s := &InventoryConfigurationService{
		InventoryConfigurationRepository: inventoryConfigurationRepository,
		MongoStorageManagerClient:        mongoStorageManager,
	}
	return s
}

func (c InventoryConfigurationService) CreateNewConfiguration(ctx context.Context, inventoryConfiguration request_dto.CreateNewConfigurationRequestBody) *dto.ErrorResponseDto {
	methodName := "AddNewInventory"
	log := logger.GetLogger()

	errorDto := c.InventoryConfigurationRepository.CreateNewConfiguration(ctx, inventoryConfiguration)
	if errorDto != nil {
		log.Error("Inside " + methodName + " error occurred when trying to create new base configuration: " + inventoryConfiguration.InventoryName)
		return errorDto
	}
	log.Info("Inside " + methodName + "successfully created base configuration: " + inventoryConfiguration.InventoryName)

	createCollectionError := c.MongoStorageManagerClient.CreateCollection(context.Background(), inventoryConfiguration.InventoryName, inventoryConfiguration.JsonSchema, inventoryConfiguration.InventoryIdentifiers)
	if createCollectionError != nil {
		log.Error("Inside " + methodName + " error occurred when trying to create base configuration collection: " + constants.InventoryCollectionNamePrefix + inventoryConfiguration.InventoryName)
		return createCollectionError
	}
	log.Info("Inside " + methodName + " successfully created new base configuration collection: " + constants.InventoryCollectionNamePrefix + inventoryConfiguration.InventoryName)

	return nil

}

func (c InventoryConfigurationService) GetInventoryConfiguration(ctx context.Context, inventoryName string) (*response_dto.InventoryConfigurationResponseDto, *dto.ErrorResponseDto) {
	methodName := "GetInventoryConfiguration"
	log := logger.GetLogger()

	log.Info("Inside "+methodName+" getting base configuration for :", inventoryName)

	inventoryConfiguration, err := c.InventoryConfigurationRepository.FetchInventoryConfigurationByName(ctx, inventoryName)
	if err != nil {
		log.Error("Inside "+methodName+" error while fetching base configuration for :", inventoryName)
		return nil, err
	}
	configurationResponse, typeErr := utils.TypeConverter[response_dto.InventoryConfigurationResponseDto](inventoryConfiguration)
	if typeErr != nil {
		var domainErr dto.ErrorResponseDto
		domainErr.SetError(status_code.IMS500)
		log.Error("Inside " + methodName + " error in type converter, json marshal/unmarshal failed: " + inventoryConfiguration.InventoryName)
		return nil, &domainErr
	}
	isDeletedErr := IsInventoryConfigurationDeleted(*inventoryConfiguration)
	return configurationResponse, isDeletedErr

}

func (c InventoryConfigurationService) GetAllInventoryConfiguration(ctx context.Context) ([]response_dto.InventoryConfigurationResponseDto, *dto.ErrorResponseDto) {
	methodName := "GetAllInventoryConfiguration"
	log := logger.GetLogger()

	log.Info("Inside " + methodName + " getting all base configuration")
	inventoryConfigurations, err := c.InventoryConfigurationRepository.FetchAllInventoryConfiguration(ctx)
	if err != nil {
		log.Error("Inside " + methodName + " error while fetching all base configuration")
	}

	//Remove deleted base configurations
	var activeInventoryConfigurations []response_dto.InventoryConfigurationResponseDto
	for _, inventoryConfiguration := range inventoryConfigurations {
		isDeletedErr := IsInventoryConfigurationDeleted(inventoryConfiguration)
		if isDeletedErr != nil {
			continue
		} else {
			configurationResponse, typeErr := utils.TypeConverter[response_dto.InventoryConfigurationResponseDto](inventoryConfiguration)
			if typeErr != nil {
				var domainErr dto.ErrorResponseDto
				domainErr.SetError(status_code.IMS500)
				log.Error("Inside " + methodName + " error in type converter, json marshal/unmarshal failed: " + inventoryConfiguration.InventoryName)
				return nil, &domainErr
			}
			activeInventoryConfigurations = append(activeInventoryConfigurations, *configurationResponse)
		}
	}
	return activeInventoryConfigurations, err
}

func (c InventoryConfigurationService) DeleteInventoryConfiguration(ctx context.Context, inventoryName string) *dto.ErrorResponseDto {
	methodName := "DeleteInventoryConfiguration"
	log := logger.GetLogger()

	log.Info("Inside "+methodName+" deleting inventory configuration for :", inventoryName)
	err := c.InventoryConfigurationRepository.DeleteInventoryConfigurationByName(ctx, inventoryName)
	if err != nil {
		log.Error("Inside "+methodName+" error while deleting inventory configuration for :", inventoryName)
	}
	return err
}

func IsInventoryConfigurationDeleted(inventoryConfiguration inventoryServiceDto.InventoryConfiguration) *dto.ErrorResponseDto {
	methodName := "IsInventoryConfigurationDeleted"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto

	if inventoryConfiguration.IsDeleted == true {
		log.Info("Inside "+methodName+" inventory configuration is deleted for inventoryName :", inventoryConfiguration.InventoryName)
		domainErr.SetError(status_code.IMS204)
		return &domainErr
	}
	return nil
}
