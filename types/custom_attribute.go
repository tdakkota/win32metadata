package types

// CustomAttribute is a II.22.10 CustomAttribute representation.
type CustomAttribute struct {
	Parent HasCustomAttribute
	Type   CustomAttributeType
	Value  Blob
}
