package types

// MemberRef is a II.22.25 MemberRef representation.
type MemberRef struct {
	Class     MemberRefParent
	Name      string
	Signature Signature
}
