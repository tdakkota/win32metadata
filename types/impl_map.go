package types

// ImplMap is a II.22.22 ImplMap representation.
type ImplMap struct {
	MappingFlags    PInvokeAttributes
	MemberForwarded MemberForwarded
	ImportName      string
	ImportScope     Index `table:"ModuleRef"`
}
