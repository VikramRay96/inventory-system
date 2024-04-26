package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/common/pkg/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
)

//go:generate mockgen -destination=mocks/mock_mongo_storage_manager.go -package=mocks . IMongoStorageManager
type IMongoStorageManager interface {
	CreateCollection(ctx context.Context, collectionString string, validation bson.M, inventoryIdentifier []request_dto.InventoryIdentifier) *dto.ErrorResponseDto
}
