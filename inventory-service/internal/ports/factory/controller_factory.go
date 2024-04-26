package factory

import (
	"inventory-system/inventory-service/internal/domain/factory"
	"inventory-system/inventory-service/internal/ports/controller"
	"inventory-system/inventory-service/internal/ports/utils"
	"sync"
)

var (
	controllerFacade     *Controllers
	controllerFacadeOnce sync.Once
)

type Controllers struct {
	HealthController                 *controller.HealthController
	InventoryController              *controller.InventoryController
	InventoryConfigurationController *controller.InventoryConfigurationController
}

// GetControllerFacade creates singleton instance of  facade controllers used by routes
func GetControllerFacade() *Controllers {
	if controllerFacade == nil {
		controllerFacadeOnce.Do(func() {
			serviceFacade := factory.GetServices()
			requestValidator := utils.NewRequestValidator()
			healthController := controller.NewHealthController()
			InventoryConfigurationController := controller.NewInventoryConfigurationController(serviceFacade.InventoryConfigurationService, requestValidator)
			InventoryController := controller.NewInventoryController(serviceFacade.InventoryService)
			controllerFacade = &Controllers{
				HealthController:                 healthController,
				InventoryConfigurationController: InventoryConfigurationController,
				InventoryController:              InventoryController,
			}
		})

	}
	return controllerFacade
}
