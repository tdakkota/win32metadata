package types

// Property is a II.22.34 Property representation.
type Property struct {
	Flags PropertyAttributes
	Name  string
	Type  Signature
}
