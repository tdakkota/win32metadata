package types

// FieldRVA is a II.22.18 FieldRVA representation.
type FieldRVA struct {
	RVA   uint32
	Field Index `table:"Field"`
}
