package types

// AssemblyRefProcessor is a II.22.7 AssemblyRefProcessor representation.
type AssemblyRefProcessor struct {
	Processor   uint32
	AssemblyRef Index `table:"AssemblyRef"`
}
