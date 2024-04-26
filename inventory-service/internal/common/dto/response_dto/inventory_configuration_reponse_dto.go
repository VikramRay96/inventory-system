package response_dto

import (
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type InventoryConfigurationResponseDto struct {
	InventoryName        string                            `bson:"inventory_name" json:"inventory_name"`
	InventoryIdentifiers []request_dto.InventoryIdentifier `bson:"inventory_identifiers" json:"inventory_identifiers"`
	CreatedBy            string                            `bson:"created_by" json:"created_by"`
	CreatedOn            time.Time                         `bson:"created_on" json:"created_on"`
	JsonSchema           bson.M                            `bson:"json_schema" json:"json_schema"`
	Pagination           bool                              `bson:"pagination" json:"pagination"`
}
