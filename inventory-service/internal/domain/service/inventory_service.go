package service

import (
	"context"
	"inventory-system/common/pkg/dto"
	commonDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//go:generate mockgen -destination=mocks/mock_inventory_service.go -package=mocks . IInventoryService

type IInventoryService interface {
	CreateNewInventory(ctx context.Context, configuration interface{}, baseConfigurationName string) *dto.ErrorResponseDto
	GetInventory(ctx context.Context, baseConfigurationName string, filterAttribute map[string][]string) (bson.M, *dto.ErrorResponseDto)
	GetInventoryV2(ctx *gin.Context, from string, to string, baseConfigurationName string, filterAttribute map[string][]string) ([]bson.M, *commonDto.PaginationResponse, *dto.ErrorResponseDto)
	RemoveItemFromInventory(ctx context.Context, RemoveInventoryItemRequest *request_dto.RemoveInventoryItem, InventoryName string) *dto.ErrorResponseDto
	RemoveSubjectTopicsByLessonNameAndSubjectId(ctx context.Context, RemoveSubjectRequest *request_dto.RemoveSubjectRequest, Type string) *dto.ErrorResponseDto
	UpdateInventoryTopic(ctx context.Context, InventoryTopicUpdateRequest *request_dto.InventoryTopicUpdateRequest, TopicId string) *dto.ErrorResponseDto
	ActivateResourceById(ctx context.Context, InventoryName string, Id string) *dto.ErrorResponseDto
	UpdateInventory(Id string, InventoryName string, UpdateRequest *interface{}) *dto.ErrorResponseDto
	CreateResource(ctx context.Context, InventoryResourceCreate request_dto.InventoryResourceCreate, topicId string, contentType string, FileExtension string) (*string, *dto.ErrorResponseDto)
	GetInventoryFilter(ctx context.Context, InventoryName string, FilterName string, filters map[string][]string) ([]interface{}, *dto.ErrorResponseDto)
}
