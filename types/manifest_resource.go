package types

// ManifestResource is a II.22.24 ManifestResource representation.
type ManifestResource struct {
	Offset         uint32
	Flags          ManifestResourceAttributes
	Name           string
	Implementation Implementation
}
