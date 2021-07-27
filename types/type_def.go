package types

// TypeDef is a II.22.37 TypeDef representation.
type TypeDef struct {
	Flags         TypeAttributes
	TypeName      string
	TypeNamespace string
	Extends       TypeDefOrRef
	FieldList     List `table:"Field"`
	MethodList    List `table:"MethodDef"`
}
