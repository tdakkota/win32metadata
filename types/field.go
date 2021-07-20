package types

type (
	// FieldAttributes describes II.23.1.5 Flags for fields [FieldAttributes].
	FieldAttributes uint16

	// Field is a II.22.15 Field representation.
	Field struct {
		Flags     FieldAttributes
		Name      string
		Signature Signature
	}
)
