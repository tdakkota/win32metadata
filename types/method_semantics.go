package types

// MethodSemantics is a II.22.28 MethodSemantics representation.
type MethodSemantics struct {
	Semantics   MethodSemanticsAttributes
	Method      Index `table:"MethodDef"`
	Association HasSemantics
}
