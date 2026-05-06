package handlers

import (
	"encoding/json"
	"log"

	"test-microservices-rabbit/events"
)

// MessageEventHandler gestisce eventi legati ai messaggi
type MessageEventHandler struct{}

// Costruttore (utile se in futuro aggiungi dipendenze)
func NewMessageEventHandler() *MessageEventHandler {
	return &MessageEventHandler{}
}

// HandleMessageSent gestisce evento "message.sent"
func (h *MessageEventHandler) HandleMessageSent(data json.RawMessage) {

	var payload events.MessageSent

	// parsing payload JSON → struct
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Println("Invalid payload:", err)
		return
	}

	// business logic
	log.Printf("Message received: %s\n", payload.Message)
}
