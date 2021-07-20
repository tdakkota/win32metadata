package types

type (
	// PropertyAttributes describes II.23.1.14 Flags for properties [PropertyAttributes].
	PropertyAttributes uint16

	// Property is a II.22.34 Property representation.
	Property struct {
		Flags PropertyAttributes
		Name  string
		Type  Signature
	}
)
