package types

// Assembly is a II.22.2 Assembly representation.
type Assembly struct {
	HashAlgId AssemblyHashAlgorithm
	Version   uint64
	Flags     AssemblyFlags
	PublicKey Blob
	Name      string
	Culture   string
}
