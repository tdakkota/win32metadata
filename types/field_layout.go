package types

// FieldLayout is a II.22.16 FieldLayout representation.
type FieldLayout struct {
	Offset uint32
	Field  Index `table:"Field"`
}
