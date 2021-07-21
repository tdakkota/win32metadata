package types

// InterfaceImpl is a II.22.23 InterfaceImpl representation.
type InterfaceImpl struct {
	Class     Index `table:"TypeDef"`
	Interface TypeDefOrRef
}
