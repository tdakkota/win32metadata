package types

// MethodSpec is a II.22.29 MethodSpec representation.
type MethodSpec struct {
	Method        MethodDefOrRef
	Instantiation Blob
}
