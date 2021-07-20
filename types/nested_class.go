package types

// NestedClass is a II.22.32 NestedClass representation.
type NestedClass struct {
	NestedClass    Index `table:"TypeDef"`
	EnclosingClass Index `table:"TypeDef"`
}
