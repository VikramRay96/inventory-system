package request_dto

import "go.mongodb.org/mongo-driver/bson"

type CreateNewConfigurationRequestBody struct {
	InventoryName        string                `json:"inventory_name" validate:"required"`
	CreatedBy            string                `json:"created_by" validate:"required"`
	JsonSchema           bson.M                `json:"json_schema" validate:"required"`
	InventoryIdentifiers []InventoryIdentifier `json:"inventory_identifiers" validate:"required"`
}

type InventoryIdentifier struct {
	Key      string `json:"key" bson:"key"`
	IsUnique bool   `json:"is_unique" bson:"is_unique"`
}
