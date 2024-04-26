package dto

import "time"

type EventRequest struct {
	ServiceName string    `json:"service_name"`
	Action      string    `json:"action"`
	EntityId    string    `json:"entity_id"`
	SubentityId string    `json:"sub_entity_id"`
	EntityType  string    `json:"entity_type"`
	Time        time.Time `json:"time"`
	Date        string    `json:"date"`
	LogLevel    string    `json:"log_level"`
}
