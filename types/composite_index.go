package types

// TypeDefOrRef represents index one of TypeDef, TypeRef, or TypeSpec tables.
type TypeDefOrRef uint32

// Tag returns TypeDefOrRef tag.
// Tag table:
//
// 	TypeDef = 0
// 	TypeRef = 1
// 	TypeSpec = 2
//
func (t TypeDefOrRef) Tag() uint32 {
	return uint32(t & ((1 << 2) - 1))
}

// TableIndex returns TypeDefOrRef index.
func (t TypeDefOrRef) TableIndex() uint32 {
	return uint32((t >> 2) - 1)
}

// HasConstant represents index one of Field, Param, or Property tables.
type HasConstant uint32

// Tag returns HasConstant tag.
// Tag table:
//
// 	Field = 0
// 	Param = 1
// 	Property = 2
//
func (t HasConstant) Tag() uint32 {
	return uint32(t & ((1 << 2) - 1))
}

// TableIndex returns HasConstant index.
func (t HasConstant) TableIndex() uint32 {
	return uint32((t >> 2) - 1)
}

// HasCustomAttribute represents index one of MethodDef, Field, TypeRef, TypeDef, Param, InterfaceImpl, MemberRef,
// Module, Permission, Property, Event, StandAloneSig, ModuleRef, TypeSpec, Assembly, AssemblyRef, File,
// ExportedType, ManifestResource, GenericParam, GenericParamConstraint, MethodSpec tables.
type HasCustomAttribute uint32

// Tag returns HasCustomAttribute tag.
// Tag table:
//
// 	MethodDef = 0
// 	Field = 1
// 	TypeRef = 2
// 	TypeDef = 3
// 	Param = 4
// 	InterfaceImpl = 5
// 	MemberRef = 6
// 	Module = 7
// 	Permission = 8
// 	Property = 9
// 	Event = 10
// 	StandAloneSig = 11
// 	ModuleRef = 12
// 	TypeSpec = 13
// 	Assembly = 14
// 	AssemblyRef = 15
// 	File = 16
// 	ExportedType = 17
// 	ManifestResource = 18
// 	GenericParam = 19
// 	GenericParamConstraint = 20
// 	MethodSpec = 21
//
func (t HasCustomAttribute) Tag() uint32 {
	return uint32(t & ((1 << 5) - 1))
}

// TableIndex returns HasCustomAttribute index.
func (t HasCustomAttribute) TableIndex() uint32 {
	return uint32((t >> 5) - 1)
}

// MemberRefParent represents index one of TypeDef, TypeRef, ModuleRef, MethodDef, TypeSpec tables.
type MemberRefParent uint32

// Tag returns MemberRefParent tag.
// Tag table:
//
// 	TypeDef 0
// 	TypeRef 1
// 	ModuleRef 2
// 	MethodDef 3
// 	TypeSpec 4
//
func (t MemberRefParent) Tag() uint32 {
	return uint32(t & ((1 << 3) - 1))
}

// TableIndex returns MemberRefParent index.
func (t MemberRefParent) TableIndex() uint32 {
	return uint32((t >> 3) - 1)
}

// MemberForwarded represents index one of Field, MethodDef tables.
type MemberForwarded uint32

// Tag returns MemberForwarded tag.
// Tag table:
//
// 	Field = 0
// 	MethodDef = 1
//
func (t MemberForwarded) Tag() uint32 {
	return uint32(t & ((1 << 1) - 1))
}

// TableIndex returns MemberForwarded index.
func (t MemberForwarded) TableIndex() uint32 {
	return uint32((t >> 1) - 1)
}

// CustomAttributeType represents index one of MethodDef, MemberRef tables.
type CustomAttributeType uint32

// Tag returns CustomAttributeType tag.
// Tag table:
//
// 	Not used = 0
// 	Not used = 1
// 	MethodDef = 2
// 	MemberRef = 3
// 	Not used = 4
//
func (t CustomAttributeType) Tag() uint32 {
	return uint32(t & ((1 << 3) - 1))
}

// TableIndex returns CustomAttributeType index.
func (t CustomAttributeType) TableIndex() uint32 {
	return uint32((t >> 3) - 1)
}

// ResolutionScope represents index one of Module, ModuleRef, AssemblyRef, TypeRef tables.
type ResolutionScope uint32

// Tag returns ResolutionScope tag.
// Tag table:
//
// 	Module = 0
// 	ModuleRef = 1
// 	AssemblyRef = 2
// 	TypeRef = 3
//
func (t ResolutionScope) Tag() uint32 {
	return uint32(t & ((1 << 2) - 1))
}

// TableIndex returns ResolutionScope index.
func (t ResolutionScope) TableIndex() uint32 {
	return uint32((t >> 2) - 1)
}

// TypeOrMethodDef represents index one of TypeDef, MethodDef tables.
type TypeOrMethodDef uint32

// Tag returns TypeOrMethodDef tag.
// Tag table:
//
// 	TypeDef = 0
// 	MethodDef = 1
//
func (t TypeOrMethodDef) Tag() uint32 {
	return uint32(t & ((1 << 1) - 1))
}

// TableIndex returns TypeOrMethodDef index.
func (t TypeOrMethodDef) TableIndex() uint32 {
	return uint32((t >> 1) - 1)
}
