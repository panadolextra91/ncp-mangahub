package models

// Event represents the canonical decoupled payload that passes through the internal Event Bus.
type Event struct {
	Topic   string
	Payload interface{}
}
