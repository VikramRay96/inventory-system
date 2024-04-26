package factories

import (
	"inventory-system/inventory-service/internal/adapters/repository"
	"inventory-system/inventory-service/internal/adapters/repository/impl"
	"sync"
)

type Repositories struct {
	InventoryConfigurationRepository repository.IInventoryConfigurationRepository
	InventoryRepository              repository.IInventoryRepository
}

var (
	repositories     *Repositories
	repositoriesOnce sync.Once
)

func GetRepositories() *Repositories {
	if repositories == nil {
		repositoriesOnce.Do(func() {
			inventoryConfigurationRepository := impl.NewInventoryConfigurationRepository()
			inventoryRepository := impl.NewInventoryRepository()
			repositories = &Repositories{
				InventoryConfigurationRepository: inventoryConfigurationRepository,
				InventoryRepository:              inventoryRepository,
			}

		})
	}
	return repositories
}
