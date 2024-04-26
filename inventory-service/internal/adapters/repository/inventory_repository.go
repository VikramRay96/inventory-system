package repository

import (
	"context"
	"inventory-system/common/pkg/dto"
	"inventory-system/inventory-service/internal/adapters/models"
	commonDto "inventory-system/inventory-service/internal/common/dto"
	"go.mongodb.org/mongo-driver/bson"
)

//go:generate mockgen -destination=mocks/mock_inventory_repository.go -package=mocks . IInventoryRepository
type IInventoryRepository interface {
	CreateNewInventoryGivenInventoryName(ctx context.Context, item interface{}, inventoryName string) *dto.ErrorResponseDto
	FetchInventory(ctx context.Context, inventoryName string, uniqueKey string, uniqueValue string) (bson.M, *dto.ErrorResponseDto)
	FetchInventoryList(ctx context.Context, from string, to string, inventoryName string, filterMap map[string][]string,pagination commonDto.Pagination) ([]bson.M, *commonDto.PaginationResponse, *dto.ErrorResponseDto)
	RemoveItemFromInventory(ctx context.Context, RemoveItemModel *models.RemoveInventoryItem, InventoryName string) *dto.ErrorResponseDto
	RemoveSubjectTopicsByLessonNameAndSubjectId(ctx context.Context, model *models.RemoveSubjectRequestModel, Type string) *dto.ErrorResponseDto
	UpdateInventoryTopic(ctx context.Context, InventoryTopicUpdateModel *models.InventoryTopicUpdateRequest, TopicId string) *dto.ErrorResponseDto
	ActivateResourceById(ctx context.Context, InventoryName string, Id string) *dto.ErrorResponseDto
	UpdateInventory(Id string, InventoryName string, UpdateRequest *interface{}) *dto.ErrorResponseDto
	GetInventoryFilter(ctx context.Context, InventoryName string, FilterName string,filters bson.M) ([]interface{}, *dto.ErrorResponseDto)
}
