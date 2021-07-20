package types

type (
	// ParamAttributes describes II.23.1.13 Flags for param [ParamAttributes].
	ParamAttributes uint16

	// Param is a II.22.33 Param representation.
	Param struct {
		Flags    FieldAttributes
		Sequence uint16
		Name     string
	}
)
