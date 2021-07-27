package types

// MethodDef is a II.22.26 MethodDef representation.
type MethodDef struct {
	RVA       uint32
	ImplFlags MethodImplAttributes
	Flags     MethodAttributes
	Name      string
	Signature Signature
	ParamList List `table:"Param"`
}
