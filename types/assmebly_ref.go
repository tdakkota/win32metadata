package types

type (
	// AssemblyFlags describes II.23.1.2 Values for AssemblyFlags.
	AssemblyFlags uint32

	// AssemblyRef is a II.22.5 AssemblyRef representation.
	AssemblyRef struct {
		MajorVersion,
		MinorVersion,
		BuildNumber,
		RevisionNumber uint16
		Flags            AssemblyFlags
		PublicKeyOrToken Blob
		Name             string
		Culture          string
		HashValue        Blob
	}
)
