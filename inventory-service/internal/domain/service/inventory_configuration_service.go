package service

import (
	"context"
	"inventory-system/common/pkg/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/dto/response_dto"
)

//go:generate mockgen -destination=mocks/mock_inventory_configuration_service.go -package=mocks . IInventoryConfigurationService
type IInventoryConfigurationService interface {
	CreateNewConfiguration(ctx context.Context, configurationDto request_dto.CreateNewConfigurationRequestBody) *dto.ErrorResponseDto
	GetInventoryConfiguration(ctx context.Context, baseConfigurationName string) (*response_dto.InventoryConfigurationResponseDto, *dto.ErrorResponseDto)
	DeleteInventoryConfiguration(ctx context.Context, baseConfigurationName string) *dto.ErrorResponseDto
	GetAllInventoryConfiguration(ctx context.Context) ([]response_dto.InventoryConfigurationResponseDto, *dto.ErrorResponseDto)
}
