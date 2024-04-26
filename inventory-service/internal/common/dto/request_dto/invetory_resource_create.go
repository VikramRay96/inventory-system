package request_dto

type InventoryResourceCreate struct {
	ServiceName   string `json:"service_name"`
	FlowType      string `json:"flow_type"`
	TopicId       string `json:"topic_id"`
	TopicName     string `json:"topic_name"`
	Resource      []byte `json:"resource"`
	FileType      string `json:"file_type"`
	CreatedBy     string `json:"created_by"`
	FileRequestId string `json:"file_request_id"`
}
