package factories

import (
	"inventory-system/inventory-service/internal/adapters/client"
	clientImpl "inventory-system/inventory-service/internal/adapters/client/impl"
	"sync"
)

type Clients struct {
	MongoStorageManagerClient client.IMongoStorageManager
	S3Client                  client.IS3
}

var (
	clients     *Clients
	clientsOnce sync.Once
)

func GetClients() *Clients {
	if clients == nil {
		clientsOnce.Do(func() {
			mongoStorageMangerClient := clientImpl.NewStorageManagerService()
			s3Client := clientImpl.NewS3()
			clients = &Clients{
				MongoStorageManagerClient: mongoStorageMangerClient,
				S3Client:                  s3Client,
			}

		})
	}
	return clients
}
