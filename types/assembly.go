package types

type (
	// AssemblyHashAlgorithm describes II.23.1.1 Values for AssemblyHashAlgorithm.
	AssemblyHashAlgorithm uint32

	// Assembly is a II.22.2 Assembly representation.
	Assembly struct {
		HashAlgId      AssemblyHashAlgorithm
		MajorVersion   uint16
		MinorVersion   uint16
		BuildNumber    uint16
		RevisionNumber uint16
		Flags          AssemblyFlags
		PublicKey      Blob
		Name           string
		Culture        string
	}
)