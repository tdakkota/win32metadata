package types

// EventMap is a II.22.12 EventMap representation.
type EventMap struct {
	Parent    Index `table:"TypeDef"`
	EventList List  `table:"Event"`
}
