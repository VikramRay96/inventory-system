package repository

import (
	"context"
	"inventory-system/common/pkg/dto"
	ConfigurationServiceDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
)

//go:generate mockgen -destination=mocks/mock_inventory_configuration_repository.go -package=mocks . IInventoryConfigurationRepository
type IInventoryConfigurationRepository interface {
	CreateNewConfiguration(ctx context.Context, baseConfiguration request_dto.CreateNewConfigurationRequestBody) *dto.ErrorResponseDto
	FetchInventoryConfigurationByName(ctx context.Context, name string) (*ConfigurationServiceDto.InventoryConfiguration, *dto.ErrorResponseDto)
	FetchAllInventoryConfiguration(ctx context.Context) ([]ConfigurationServiceDto.InventoryConfiguration, *dto.ErrorResponseDto)
	DeleteInventoryConfigurationByName(ctx context.Context, name string) *dto.ErrorResponseDto
}
