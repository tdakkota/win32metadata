package types

type (
	// MethodSemanticsAttributes describes II.23.1.12 Flags for MethodSemantics [MethodSemanticsAttributes].
	MethodSemanticsAttributes uint16

	// MethodSemantics is a II.22.28 MethodSemantics representation.
	MethodSemantics struct {
		Semantics   MethodSemanticsAttributes
		Method      Index `table:"MethodDef"`
		Association HasSemantics
	}
)
