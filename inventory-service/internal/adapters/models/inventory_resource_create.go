package models

import (
	"time"
)

type InventoryResourceCreateModel struct {
	ServiceName string    `json:"service_name" bson:"service_name"`
	FlowType    string    `json:"flow_type" bson:"flow_type"`
	TopicId     string    `json:"topic_id" bson:"topic_id"`
	FileType    string    `json:"file_type" bson:"file_type"`
	VideosUrl   string    `json:"videos_url" bson:"videos_url"`
	CreatedBy   string    `json:"created_by" bson:"created_by"`
	CreatedOn   time.Time `bson:"created_on"`
}
