package types

// ExportedType is a II.22.14 ExportedType representation.
type ExportedType struct {
	Flags          TypeAttributes
	TypeDefId      uint32
	TypeName       string
	TypeNamespace  string
	Implementation Implementation
}
