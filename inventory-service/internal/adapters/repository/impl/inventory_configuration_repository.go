package impl

import (
	"context"
	"inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/common/pkg/utils"
	"inventory-system/inventory-service/internal/adapters/db"
	"inventory-system/inventory-service/internal/adapters/models"
	ConfigurationServiceDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryConfigurationRepository struct {
}

func NewInventoryConfigurationRepository() *InventoryConfigurationRepository {
	repo := &InventoryConfigurationRepository{}
	return repo
}
func (c InventoryConfigurationRepository) CreateNewConfiguration(ctx context.Context, baseConfiguration request_dto.CreateNewConfigurationRequestBody) *dto.ErrorResponseDto {
	methodName := "CreateNewConfiguration"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto

	configurationModel, err := utils.TypeConverter[models.InventoryConfiguration](baseConfiguration)
	if err != nil {
		adapterErr.SetError(status_code.IMS500)
		log.Error("Inside " + methodName + " error when converting dto to model, json marshal/unmarshal failed: " + baseConfiguration.InventoryName)
		return &adapterErr
	}
	createdOn := time.Now()
	configurationModel.CreatedOn = createdOn

	_, err = db.GetDb().Collection(constants.InventoryConfigurationCollectionName).InsertOne(ctx, configurationModel)
	if err != nil {
		errorData := err.(mongo.ServerError)
		if errorData.HasErrorCode(constants.MongoDuplicateEntryErrorCode) {
			log.Error("Inside " + methodName + " error duplicate entry occurred when trying to create inventory configuration: " + baseConfiguration.InventoryName)
			adapterErr.SetError(status_code.IMS104)
			return &adapterErr
		}

		log.Error("Inside "+methodName+" error: ", err.Error(), " while creating the inventory configuration: "+baseConfiguration.InventoryName)
		adapterErr.SetError(status_code.IMS102)
		return &adapterErr
	}

	return nil
}
func (c InventoryConfigurationRepository) FetchInventoryConfigurationByName(ctx context.Context, inventoryName string) (*ConfigurationServiceDto.InventoryConfiguration, *dto.ErrorResponseDto) {
	methodName := "FetchInventoryConfigurationByName"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto

	filter := bson.M{"inventory_name": inventoryName}
	var inventoryConfiguration ConfigurationServiceDto.InventoryConfiguration
	err := db.GetDb().Collection(constants.InventoryConfigurationCollectionName).FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&inventoryConfiguration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("Inside "+methodName+" no documents exists in base configuration for filter : ", filter)
			adapterErr.SetError(status_code.IMS114)
			return nil, &adapterErr
		}
		log.Info("Inside "+methodName+" error while fetching inventory configuration for filter : ", filter, " error : ", err.Error())
		adapterErr.SetError(status_code.IMS105)
		return nil, &adapterErr
	}
	return &inventoryConfiguration, nil
}
func (c InventoryConfigurationRepository) FetchAllInventoryConfiguration(ctx context.Context) ([]ConfigurationServiceDto.InventoryConfiguration, *dto.ErrorResponseDto) {
	methodName := "FetchAllInventoryConfiguration"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto

	var inventoryConfigurations []ConfigurationServiceDto.InventoryConfiguration

	cur, err := db.GetDb().Collection(constants.InventoryConfigurationCollectionName).Find(ctx, bson.D{}, options.Find().SetProjection(bson.M{"_id": 0}))
	if err != nil {
		log.Error("Inside " + methodName + " error while fetching all inventory configurations")
		adapterErr.SetError(status_code.IMS106)
		return nil, &adapterErr

	}

	err = cur.All(ctx, &inventoryConfigurations)
	if err != nil {
		log.Error("Inside " + methodName + " error while decoding all base configurations")
		adapterErr.SetError(status_code.IMS106)
		return nil, &adapterErr

	}

	return inventoryConfigurations, nil
}
func (c InventoryConfigurationRepository) DeleteInventoryConfigurationByName(ctx context.Context, inventoryName string) *dto.ErrorResponseDto {
	methodName := "DeleteInventoryConfigurationByName"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto

	filter := bson.D{{"inventory_name", inventoryName}}

	baseConfiguration, fetchErr := c.FetchInventoryConfigurationByName(ctx, inventoryName)
	if fetchErr != nil {
		log.Error("Inside ", methodName, " error while fetching inventory configuration for inventoryName: ", inventoryName)
		return fetchErr
	}
	if baseConfiguration.IsDeleted == true {
		log.Info("Inside ", methodName, "inventory configuration already deleted inventoryName: ", inventoryName)
		adapterErr.SetError(status_code.IMS205)
		return &adapterErr
	}

	updateBody := bson.D{
		{"$set", bson.M{"is_deleted": true}},
	}
	resp, err := db.GetDb().Collection(constants.InventoryConfigurationCollectionName).UpdateOne(ctx, filter, updateBody)
	if resp.MatchedCount == 0 {
		log.Info("Inside ", methodName, " no documents found for inventory configuration for inventoryName: ", inventoryName)
		adapterErr.SetError(status_code.IMS114)
		return &adapterErr
	}

	if err != nil {
		log.Info("Inside ", methodName, " error occurred while deleting for inventory configuration for inventoryName: ", inventoryName)
		adapterErr.SetError(status_code.IMS107)
		return &adapterErr
	}
	log.Info("Inside ", methodName, " success while deleting inventory configuration for inventoryName: ", inventoryName)
	return nil
}
