package events

import (
	"encoding/json"
)

type Event struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Version       string          `json:"version"`
	Source        string          `json:"source"`
	Timestamp     int64           `json:"timestamp"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Data          json.RawMessage `json:"data"`
}

type MessageSent struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
