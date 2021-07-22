package types

type (
	// TypeAttributes describes II.23.1.15 Flags for types [TypeAttributes].
	TypeAttributes uint32

	// TypeDef is a II.22.37 TypeDef representation.
	TypeDef struct {
		Flags         TypeAttributes
		TypeName      string
		TypeNamespace string
		Extends       TypeDefOrRef
		FieldList     List `table:"Field"`
		MethodList    List `table:"MethodDef"`
	}
)
