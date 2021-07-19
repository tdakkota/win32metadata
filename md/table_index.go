package md

func (h *TablesHeader) computeIndexes() {
	// Ported from https://github.com/microsoft/windows-rs/blob/master/crates/gen/src/parser/file.rs#L424-L729
	var (
		empty           Table
		blobIndexSize   = h.BlobIndexSize()
		stringIndexSize = h.StringIndexSize()
		guidIndexSize   = h.GUIDIndexSize()
	)

	typeDefOrRef := compositeIndexSize(
		h.Tables[TypeDef],
		h.Tables[TypeRef],
		h.Tables[TypeSpec],
	)
	hasConstant := compositeIndexSize(
		h.Tables[TypeDef],
		h.Tables[TypeRef],
		h.Tables[Property],
	)
	hasCustomAttribute := compositeIndexSize(
		h.Tables[MethodDef],
		h.Tables[Field],
		h.Tables[TypeRef],
		h.Tables[TypeDef],
		h.Tables[Param],
		h.Tables[InterfaceImpl],
		h.Tables[MemberRef],
		h.Tables[Module],
		h.Tables[Property],
		h.Tables[Event],
		h.Tables[StandaloneSig],
		h.Tables[ModuleRef],
		h.Tables[TypeSpec],
		h.Tables[Assembly],
		h.Tables[AssemblyRef],
		h.Tables[File],
		h.Tables[ExportedType],
		h.Tables[ManifestResource],
		h.Tables[GenericParam],
		h.Tables[GenericParamConstraint],
		h.Tables[MethodSpec],
	)
	hasFieldMarshal := compositeIndexSize(
		h.Tables[Field],
		h.Tables[Param],
	)
	hasDeclSecurity := compositeIndexSize(
		h.Tables[TypeDef],
		h.Tables[MethodDef],
		h.Tables[Assembly],
	)
	memberRefParent := compositeIndexSize(
		h.Tables[TypeDef],
		h.Tables[TypeRef],
		h.Tables[ModuleRef],
		h.Tables[MethodDef],
		h.Tables[TypeSpec],
	)
	hasSemantics := compositeIndexSize(
		h.Tables[Event],
		h.Tables[Property],
	)
	methodDefOrRef := compositeIndexSize(
		h.Tables[MethodDef],
		h.Tables[MemberRef],
	)
	memberForwarded := compositeIndexSize(
		h.Tables[Field],
		h.Tables[MethodDef],
	)
	implementation := compositeIndexSize(
		h.Tables[File],
		h.Tables[AssemblyRef],
		h.Tables[ExportedType],
	)
	customAttributeType := compositeIndexSize(
		h.Tables[MethodDef],
		h.Tables[MemberRef],
		empty,
		empty,
		empty,
	)
	resolutionScope := compositeIndexSize(
		h.Tables[Module],
		h.Tables[ModuleRef],
		h.Tables[AssemblyRef],
		h.Tables[TypeRef],
	)
	typeOrMethodDef := compositeIndexSize(
		h.Tables[TypeDef],
		h.Tables[MethodDef],
	)

	h.Tables[Assembly].SetRowType([6]uint32{4, 8, 4, blobIndexSize, stringIndexSize, stringIndexSize})
	h.Tables[AssemblyOs].SetRowType([6]uint32{4, 4, 4})
	h.Tables[AssemblyProcessor].SetRowType([6]uint32{4})
	h.Tables[AssemblyRef].SetRowType([6]uint32{8, 4, blobIndexSize, stringIndexSize, stringIndexSize, blobIndexSize})
	h.Tables[AssemblyRefOs].SetRowType([6]uint32{4, 4, 4, h.Tables[AssemblyRef].IndexSize()})
	h.Tables[AssemblyRefProcessor].SetRowType([6]uint32{4, h.Tables[AssemblyRef].IndexSize()})
	h.Tables[ClassLayout].SetRowType([6]uint32{2, 4, h.Tables[TypeDef].IndexSize()})
	h.Tables[Constant].SetRowType([6]uint32{2, hasConstant, blobIndexSize})
	h.Tables[CustomAttribute].SetRowType([6]uint32{hasCustomAttribute, customAttributeType, blobIndexSize})
	h.Tables[DeclSecurity].SetRowType([6]uint32{2, hasDeclSecurity, blobIndexSize})
	h.Tables[EventMap].SetRowType([6]uint32{h.Tables[TypeDef].IndexSize(), h.Tables[Event].IndexSize()})
	h.Tables[Event].SetRowType([6]uint32{2, stringIndexSize, typeDefOrRef})
	h.Tables[ExportedType].SetRowType([6]uint32{4, 4, stringIndexSize, stringIndexSize, implementation})
	h.Tables[Field].SetRowType([6]uint32{2, stringIndexSize, blobIndexSize})
	h.Tables[FieldLayout].SetRowType([6]uint32{4, h.Tables[Field].IndexSize()})
	h.Tables[FieldMarshal].SetRowType([6]uint32{hasFieldMarshal, blobIndexSize})
	h.Tables[FieldRva].SetRowType([6]uint32{4, h.Tables[Field].IndexSize()})
	h.Tables[File].SetRowType([6]uint32{4, stringIndexSize, blobIndexSize})
	h.Tables[GenericParam].SetRowType([6]uint32{2, 2, typeOrMethodDef, stringIndexSize})
	h.Tables[GenericParamConstraint].SetRowType([6]uint32{h.Tables[GenericParam].IndexSize(), typeDefOrRef})
	h.Tables[ImplMap].SetRowType([6]uint32{2, memberForwarded, stringIndexSize, h.Tables[ModuleRef].IndexSize()})
	h.Tables[InterfaceImpl].SetRowType([6]uint32{h.Tables[TypeDef].IndexSize(), typeDefOrRef})
	h.Tables[ManifestResource].SetRowType([6]uint32{4, 4, stringIndexSize, implementation})
	h.Tables[MemberRef].SetRowType([6]uint32{memberRefParent, stringIndexSize, blobIndexSize})
	h.Tables[MethodDef].SetRowType([6]uint32{4, 2, 2, stringIndexSize, blobIndexSize, h.Tables[Param].IndexSize()})
	h.Tables[MethodImpl].SetRowType([6]uint32{h.Tables[TypeDef].IndexSize(), methodDefOrRef, methodDefOrRef})
	h.Tables[MethodSemantics].SetRowType([6]uint32{2, h.Tables[MethodDef].IndexSize(), hasSemantics})
	h.Tables[MethodSpec].SetRowType([6]uint32{methodDefOrRef, blobIndexSize})
	h.Tables[Module].SetRowType([6]uint32{2, stringIndexSize, guidIndexSize, guidIndexSize, guidIndexSize})
	h.Tables[ModuleRef].SetRowType([6]uint32{stringIndexSize})
	h.Tables[NestedClass].SetRowType([6]uint32{h.Tables[TypeDef].IndexSize(), h.Tables[TypeDef].IndexSize()})
	h.Tables[Param].SetRowType([6]uint32{2, 2, stringIndexSize})
	h.Tables[Property].SetRowType([6]uint32{2, stringIndexSize, blobIndexSize})
	h.Tables[PropertyMap].SetRowType([6]uint32{h.Tables[TypeDef].IndexSize(), h.Tables[Property].IndexSize()})
	h.Tables[StandaloneSig].SetRowType([6]uint32{blobIndexSize})
	h.Tables[TypeDef].SetRowType([6]uint32{
		4,
		stringIndexSize,
		stringIndexSize,
		typeDefOrRef,
		h.Tables[Field].IndexSize(),
		h.Tables[MethodDef].IndexSize(),
	})
	h.Tables[TypeRef].SetRowType([6]uint32{resolutionScope, stringIndexSize, stringIndexSize})
	h.Tables[TypeSpec].SetRowType([6]uint32{blobIndexSize})
}

func compositeIndexSize(t ...Table) uint32 {
	small := func(rowCount uint32, bits uint8) bool {
		return uint64(rowCount) < (uint64(1) << (16 - bits))
	}

	bitsNeeded := func(value int) (bits uint8) {
		value -= 1
		bits = 1
		for {
			value >>= 1
			if value == 0 {
				break
			}
			bits += 1
		}
		return bits
	}

	bits := bitsNeeded(len(t))
	for i := range t {
		if !small(t[i].RowCount, bits) {
			return 4
		}
	}
	return 2
}
