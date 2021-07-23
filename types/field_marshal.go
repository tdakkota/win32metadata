package types

// FieldMarshal is a II.22.17 FieldMarshal representation.
type FieldMarshal struct {
	Parent     Index // Should be HasFieldMarshal, but ECMA-335 6th-edition does not define it :(.
	NativeType Blob
}
