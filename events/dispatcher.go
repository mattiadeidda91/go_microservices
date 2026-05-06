package events

import (
	"encoding/json"
	"log"
)

// EventHandler è una funzione che gestisce un evento
// riceve il payload raw (JSON)
type EventHandler func(data json.RawMessage)

// Dispatcher gestisce la mappa eventType → handler
type Dispatcher struct {
	handlers map[string]EventHandler
}

// Crea nuovo dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]EventHandler),
	}
}

// Register associa un tipo evento a un handler
func (d *Dispatcher) Register(eventType string, handler EventHandler) {
	d.handlers[eventType] = handler
}

// Dispatch riceve un evento e chiama l'handler corretto
func (d *Dispatcher) Dispatch(event Event) {

	handler, exists := d.handlers[event.Type]

	if !exists {
		log.Println("No handler for event:", event.Type)
		return
	}

	// chiama handler passando il payload
	handler(event.Data)
}
