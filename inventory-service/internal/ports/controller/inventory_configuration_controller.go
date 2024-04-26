package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	ConfigurationServiceDto "inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"inventory-system/inventory-service/internal/domain/service"
	"inventory-system/inventory-service/internal/ports/utils"
	"net/http"
)

type InventoryConfigurationController struct {
	RequestValidator              utils.IRequestValidator
	InventoryConfigurationService service.IInventoryConfigurationService
}

func NewInventoryConfigurationController(inventoryConfigurationService service.IInventoryConfigurationService, requestValidator utils.IRequestValidator) *InventoryConfigurationController {
	controller := &InventoryConfigurationController{
		InventoryConfigurationService: inventoryConfigurationService,
		RequestValidator:              requestValidator,
	}
	return controller
}

// CreateNewConfiguration  godoc
// @Summary Create new inventory configuration
// @Description Create new inventory configuration
// @Tags InventoryConfiguration
// @Accept  json
// @Produce  json
// @Param requestBody body request_dto.CreateNewConfigurationRequestBody true "Create Inventory Configuration Request"
// @Success 200 {object} dto.ResponseDto
// @Router /inventory-service/api/v1/inventory/configurations [POST]
// CreateNewConfiguration : This function will Create New Inventory Configuration
func (bc InventoryConfigurationController) CreateNewConfiguration() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "AddNewInventory"
		log := logger.GetLogger()
		ctx := context.Background()
		var portErr dto.ErrorResponseDto

		var createNewConfigurationRequestBody ConfigurationServiceDto.CreateNewConfigurationRequestBody

		err := c.ShouldBindJSON(&createNewConfigurationRequestBody)
		if err != nil {
			log.Error("Inside "+methodName+" error while binding json Error: ", err.Error())
			portErr.SetError(status_code.IMS400)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: portErr.StatusCode,
				Message:    portErr.Message,
			})
			return
		}

		errorDto, _ := bc.RequestValidator.ValidateCreateConfigurationRequest(createNewConfigurationRequestBody)
		if errorDto != nil {
			log.Error("Inside " + methodName + " error while validating request for Create Configuration")
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}

		errorDto = bc.InventoryConfigurationService.CreateNewConfiguration(ctx, createNewConfigurationRequestBody)
		if errorDto != nil {
			log.Error("Inside "+methodName+" error while creating new configuration: ", createNewConfigurationRequestBody.InventoryName)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
	return fn
}

// GetConfiguration  godoc
// @Summary Fetch Inventory Configuration
// @Description fetch inventory configuration
// @Tags InventoryConfiguration
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Param inventoryName path string true "Inventory Key"
// @Router /inventory-service/api/v1/inventory/configurations/{inventoryName} [GET]
// GetConfiguration : This function will fetch inventory Configuration
func (bc InventoryConfigurationController) GetConfiguration() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "GetInventory"
		log := logger.GetLogger()
		ctx := context.Background()
		configurationName := c.Param("inventoryName")

		inventoryConfiguration, errDto := bc.InventoryConfigurationService.GetInventoryConfiguration(ctx, configurationName)
		if errDto != nil {
			log.Info("Inside "+methodName+" unable to fetch inventory configuration for configurationName :", configurationName)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			Data:       inventoryConfiguration,
		})
	}
	return fn
}

// GetAllConfiguration  godoc
// @Summary Fetch all Inventory configurations
// @Description Fetch all Inventory configurations
// @Tags InventoryConfiguration
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Router /inventory-service/api/v1/inventory/configurations [GET]
// GetAllConfiguration : This function will fetch all inventory Configuration
func (bc InventoryConfigurationController) GetAllConfiguration() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "GetInventory"
		log := logger.GetLogger()
		ctx := context.Background()

		inventoryConfigurations, errDto := bc.InventoryConfigurationService.GetAllInventoryConfiguration(ctx)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to fetch all base configuration")
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			Data:       inventoryConfigurations,
		})
	}
	return fn
}

// DeleteConfiguration  godoc
// @Summary Create delete inventory configuration
// @Description Create delete inventory configuration
// @Tags InventoryConfiguration
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Router /inventory-service/api/v1/inventory/configurations/{configurationName} [DELETE]
// @Param configurationName path string true "Configuration Key"
// DeleteConfiguration : This function will delete  inventory Configuration
func (bc InventoryConfigurationController) DeleteConfiguration() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "GetInventory"
		log := logger.GetLogger()
		ctx := context.Background()
		inventoryName := c.Param("inventoryName")

		errDto := bc.InventoryConfigurationService.DeleteInventoryConfiguration(ctx, inventoryName)
		if errDto != nil {
			log.Info("Inside "+methodName+" unable to fetch inventory configuration for inventoryName :", inventoryName)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
	return fn
}
