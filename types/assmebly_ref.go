package types

// AssemblyRef is a II.22.5 AssemblyRef representation.
type AssemblyRef struct {
	Version          uint64
	Flags            AssemblyFlags
	PublicKeyOrToken Blob
	Name             string
	Culture          string
	HashValue        Blob
}
