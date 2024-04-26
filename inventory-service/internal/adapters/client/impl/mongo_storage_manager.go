package impl

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/inventory-service/internal/adapters/db"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
)

type MongoStorageManager struct {
}

func NewStorageManagerService() *MongoStorageManager {
	service := &MongoStorageManager{}
	return service
}

func (s MongoStorageManager) CreateCollection(ctx context.Context, collectionString string, validation bson.M, inventoryIdentifiers []request_dto.InventoryIdentifier) *dto.ErrorResponseDto {
	methodName := "CreateCollection"
	log := logger.GetLogger()
	var ClientErr dto.ErrorResponseDto

	collectionName := constants.InventoryCollectionNamePrefix + collectionString

	collectionValidatorOptions := options.CreateCollection()
	collectionValidatorOptions.SetValidator(validation)
	collectionValidatorOptions.SetValidationLevel("strict")
	collectionValidatorOptions.SetValidationAction("error")

	err := db.GetDb().CreateCollection(ctx, collectionName, collectionValidatorOptions)
	if err != nil {
		log.Error("Inside " + methodName + " error: " + err.Error() + " occurred while creating collection: " + collectionName)
		ClientErr.SetError(status_code.IMS102)
		return &ClientErr
	}

	for _, inventoryIdentifier := range inventoryIdentifiers {
		_, err = db.GetDb().Collection(collectionName).Indexes().CreateOne(ctx, GetInventoryIndexes(inventoryIdentifier))
		if err != nil {
			log.Error("Inside " + methodName + " error: " + err.Error() + " occurred while creating index for collectionName: " + collectionName)
			ClientErr.SetError(status_code.IMS103)
			return &ClientErr
		}
	}

	return nil
}

func GetInventoryIndexes(inventoryIdentifier request_dto.InventoryIdentifier) mongo.IndexModel {
	indexName := "idx_" + inventoryIdentifier.Key
	isUnique := inventoryIdentifier.IsUnique
	indexOptions := options.IndexOptions{
		Unique: &isUnique,
		Name:   &indexName,
	}
	return mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   inventoryIdentifier.Key,
				Value: 1,
			},
		},
		Options: &indexOptions,
	}

}
