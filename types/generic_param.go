package types

type (
	// GenericParamAttributes describes II.23.1.7 Flags for Generic Parameters [GenericParamAttributes].
	GenericParamAttributes uint16

	// GenericParam is a II.22.20 GenericParam representation.
	GenericParam struct {
		Number uint16
		Flags  GenericParamAttributes
		Owner  TypeOrMethodDef
		Name   string
	}
)
