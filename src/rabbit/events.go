package rabbit

import "encoding/json"

type SystemEvent struct {
	IncludedUsers []int  `json:"included_users"`
	EventType     string `json:"event_type"`
	Data          string `json:"data"`
}

func NewSystemEvent(eventType string, includedUsers []int, data interface{}) (*SystemEvent, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &SystemEvent{IncludedUsers: includedUsers, EventType: eventType, Data: string(json_data)}, nil
}
