package types

// Field is a II.22.15 Field representation.
type Field struct {
	Flags     FieldAttributes
	Name      string
	Signature Signature
}
