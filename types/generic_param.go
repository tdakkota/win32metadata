package types

// GenericParam is a II.22.20 GenericParam representation.
type GenericParam struct {
	Number uint16
	Flags  GenericParamAttributes
	Owner  TypeOrMethodDef
	Name   string
}
