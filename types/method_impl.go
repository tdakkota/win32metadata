package types

// MethodImpl is a II.22.27 MethodImpl representation.
type MethodImpl struct {
	Class             Index `table:"TypeDef"`
	MethodBody        MethodDefOrRef
	MethodDeclaration MethodDefOrRef
}
