package models

type InventoryTopicUpdateRequest struct {
	TopicName   string         `json:"topic_name" bson:"topic_name"`
	Resources   []*interface{} `json:"resources" bson:"resources"`
	Description string         `json:"description" bson:"description"`
	Assessments []*interface{} `json:"assessments" bson:"assessments"`
}
