package types

// Event is a II.22.13 Event representation.
type Event struct {
	EventFlags EventAttributes
	Name       string
	EventType  TypeDefOrRef
}
