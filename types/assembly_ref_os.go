package types

// AssemblyRefOS is a II.22.6 AssemblyRefOS representation.
type AssemblyRefOS struct {
	OSPlatformID   uint32
	OSMajorVersion uint32
	OSMinorVersion uint32
	AssemblyRef    Index `table:"AssemblyRef"`
}
