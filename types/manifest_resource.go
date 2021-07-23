package types

type (
	// ManifestResourceAttributes describes II.23.1.9 Flags for ManifestResource [ManifestResourceAttributes].
	ManifestResourceAttributes uint32

	// ManifestResource is a II.22.24 ManifestResource representation.
	ManifestResource struct {
		Offset         uint32
		Flags          ManifestResourceAttributes
		Name           string
		Implementation Implementation
	}
)
