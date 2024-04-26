package models

type RemoveInventoryItem struct {
	ItemId string `json:"item_id" bson:"item_id"`
}
