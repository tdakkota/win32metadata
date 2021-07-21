package types

// DeclSecurity is a II.22.11 DeclSecurity representation.
type DeclSecurity struct {
	Action        uint16
	Parent        HasDeclSecurity
	PermissionSet Blob
}
