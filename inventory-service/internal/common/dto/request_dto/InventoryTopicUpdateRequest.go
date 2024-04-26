package request_dto

type InventoryTopicUpdateRequest struct {
	TopicName   string         `json:"topic_name"`
	Resources   []*interface{} `json:"resources"`
	Description string         `json:"description"`
	Assessments []*interface{} `json:"assessments"`
}
