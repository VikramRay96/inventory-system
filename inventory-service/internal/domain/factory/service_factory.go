package factory

import (
	"inventory-system/inventory-service/internal/adapters/client"
	clientImpl "inventory-system/inventory-service/internal/adapters/client/impl"
	"inventory-system/inventory-service/internal/adapters/factories"
	"inventory-system/inventory-service/internal/domain/service"
	"inventory-system/inventory-service/internal/domain/service/impl"
	"sync"
)

type Services struct {
	InventoryConfigurationService service.IInventoryConfigurationService
	MongoStorageManagerClient     client.IMongoStorageManager
	InventoryService              service.IInventoryService
	S3                            client.IS3
}

var (
	services     *Services
	servicesOnce sync.Once
)

func GetServices() *Services {
	if services == nil {
	}
	servicesOnce.Do(func() {
		repoFacade := factories.GetRepositories()
		mongoStorageManagerClient := clientImpl.NewStorageManagerService()
		s3 := clientImpl.NewS3()
		inventoryConfigurationService := impl.NewInventoryConfigurationService(repoFacade.InventoryConfigurationRepository, mongoStorageManagerClient)
		inventoryService := impl.NewInventoryService(repoFacade.InventoryRepository, inventoryConfigurationService, s3)

		services = &Services{
			InventoryConfigurationService: inventoryConfigurationService,
			MongoStorageManagerClient:     mongoStorageManagerClient,
			InventoryService:              inventoryService,
			S3:                            s3,
		}

	})
	return services
}
