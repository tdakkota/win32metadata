package types

type (
	// EventAttributes describes II.23.1.4 Flags for events [EventAttributes].
	EventAttributes uint16

	// Event is a II.22.13 Event representation.
	Event struct {
		EventFlags EventAttributes
		Name       string
		EventType  TypeDefOrRef
	}
)
