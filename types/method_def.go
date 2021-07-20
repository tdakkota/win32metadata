package types

type (
	// MethodAttributes describes II.23.1.10 Flags for methods [MethodAttributes].
	MethodAttributes uint16
	// MethodImplAttributes describes II.23.1.11 Flags for methods [MethodImplAttributes].
	MethodImplAttributes uint16

	// MethodDef is a II.22.26 MethodDef representation.
	MethodDef struct {
		RVA       uint32
		ImplFlags MethodImplAttributes
		Flags     MethodAttributes
		Name      string
		Signature Signature
		ParamList List `table:"Param"`
	}
)
