package types

// TypeRef is a II.22.38 TypeRef representation.
type TypeRef struct {
	ResolutionScope ResolutionScope
	TypeName        string
	TypeNamespace   string
}
