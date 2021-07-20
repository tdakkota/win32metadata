package types

type (
	// PInvokeAttributes describes II.23.1.8 Flags for ImplMap [PInvokeAttributes].
	PInvokeAttributes uint16

	// ImplMap is a II.22.22 ImplMap representation.
	ImplMap struct {
		MappingFlags    PInvokeAttributes
		MemberForwarded MemberForwarded
		ImportName      string
		ImportScope     Index `table:"ModuleRef"`
	}
)
