package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"time"
)

type InventoryConfiguration struct {
	InventoryName        string                            `bson:"inventory_name" json:"inventory_name"`
	CreatedBy            string                            `bson:"created_by" json:"created_by"`
	CreatedOn            time.Time                         `bson:"created_on" json:"created_on"`
	JsonSchema           bson.M                            `bson:"json_schema" json:"json_schema"`
	InventoryIdentifiers []request_dto.InventoryIdentifier `bson:"inventory_identifiers" json:"inventory_identifiers"`
	IsDeleted            bool                              `bson:"is_deleted" json:"is_deleted"`
}
