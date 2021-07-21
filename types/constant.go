package types

// Constant is a II.22.9 Constant representation.
type Constant struct {
	Type   ElementTypeKind
	Parent HasConstant
	Value  Blob
}
