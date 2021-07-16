package md

import (
	"debug/pe"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	a := require.New(t)

	f, err := pe.Open(`./testdata/.windows/winmd/Windows.Win32.winmd`)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	file, err := ParseMetadata(f)
	a.NoError(err)

	section, err := file.StreamByName("#~")
	a.NoError(err)

	var header TablesHeader
	if err := header.Decode(section); err != nil {
		t.Fatal(err)
	}

	a.Equal([4]byte{}, header.Reserved)
	a.Equal(byte(2), header.MajorVersion)
	a.Equal(byte(0), header.MinorVersion)

	var (
		tables          [GenericParamConstraint + 1]Table
		empty           Table
		blobIndexSize   = header.BlobIndexSize()
		stringIndexSize = header.StringIndexSize()
		guidIndexSize   = header.GUIDIndexSize()
	)
	for _, row := range header.Rows {
		tables[row.Tag] = Table{
			Type:     row.Tag,
			RowCount: row.Count,
		}
	}

	typeDefOrRef := compositeIndexSize(
		tables[TypeDef],
		tables[TypeRef],
		tables[TypeSpec],
	)
	hasConstant := compositeIndexSize(
		tables[TypeDef],
		tables[TypeRef],
		tables[Property],
	)
	hasCustomAttribute := compositeIndexSize(
		tables[MethodDef],
		tables[Field],
		tables[TypeRef],
		tables[TypeDef],
		tables[Param],
		tables[InterfaceImpl],
		tables[MemberRef],
		tables[Module],
		tables[Property],
		tables[Event],
		tables[StandaloneSig],
		tables[ModuleRef],
		tables[TypeSpec],
		tables[Assembly],
		tables[AssemblyRef],
		tables[File],
		tables[ExportedType],
		tables[ManifestResource],
		tables[GenericParam],
		tables[GenericParamConstraint],
		tables[MethodSpec],
	)
	hasFieldMarshal := compositeIndexSize(
		tables[Field],
		tables[Param],
	)
	hasDeclSecurity := compositeIndexSize(
		tables[TypeDef],
		tables[MethodDef],
		tables[Assembly],
	)
	memberRefParent := compositeIndexSize(
		tables[TypeDef],
		tables[TypeRef],
		tables[ModuleRef],
		tables[MethodDef],
		tables[TypeSpec],
	)
	hasSemantics := compositeIndexSize(
		tables[Event],
		tables[Property],
	)
	methodDefOrRef := compositeIndexSize(
		tables[MethodDef],
		tables[MemberRef],
	)
	memberForwarded := compositeIndexSize(
		tables[Field],
		tables[MethodDef],
	)
	implementation := compositeIndexSize(
		tables[File],
		tables[AssemblyRef],
		tables[ExportedType],
	)
	customAttributeType := compositeIndexSize(
		tables[MethodDef],
		tables[MemberRef],
		empty,
		empty,
		empty,
	)
	resolutionScope := compositeIndexSize(
		tables[Module],
		tables[ModuleRef],
		tables[AssemblyRef],
		tables[TypeRef],
	)
	typeOrMethodDef := compositeIndexSize(
		tables[TypeDef],
		tables[MethodDef],
	)

	tables[Assembly].SetRowType([6]uint32{4, 8, 4, blobIndexSize, stringIndexSize, stringIndexSize})
	tables[AssemblyOs].SetRowType([6]uint32{4, 4, 4})
	tables[AssemblyProcessor].SetRowType([6]uint32{4})
	tables[AssemblyRef].SetRowType([6]uint32{8, 4, blobIndexSize, stringIndexSize, stringIndexSize, blobIndexSize})
	tables[AssemblyRefOs].SetRowType([6]uint32{4, 4, 4, tables[AssemblyRef].IndexSize()})
	tables[AssemblyRefProcessor].SetRowType([6]uint32{4, tables[AssemblyRef].IndexSize()})
	tables[ClassLayout].SetRowType([6]uint32{2, 4, tables[TypeDef].IndexSize()})
	tables[Constant].SetRowType([6]uint32{2, hasConstant, blobIndexSize})
	tables[CustomAttribute].SetRowType([6]uint32{hasCustomAttribute, customAttributeType, blobIndexSize})
	tables[DeclSecurity].SetRowType([6]uint32{2, hasDeclSecurity, blobIndexSize})
	tables[EventMap].SetRowType([6]uint32{tables[TypeDef].IndexSize(), tables[Event].IndexSize()})
	tables[Event].SetRowType([6]uint32{2, stringIndexSize, typeDefOrRef})
	tables[ExportedType].SetRowType([6]uint32{4, 4, stringIndexSize, stringIndexSize, implementation})
	tables[Field].SetRowType([6]uint32{2, stringIndexSize, blobIndexSize})
	tables[FieldLayout].SetRowType([6]uint32{4, tables[Field].IndexSize()})
	tables[FieldMarshal].SetRowType([6]uint32{hasFieldMarshal, blobIndexSize})
	tables[FieldRva].SetRowType([6]uint32{4, tables[Field].IndexSize()})
	tables[File].SetRowType([6]uint32{4, stringIndexSize, blobIndexSize})
	tables[GenericParam].SetRowType([6]uint32{2, 2, typeOrMethodDef, stringIndexSize})
	tables[GenericParamConstraint].SetRowType([6]uint32{tables[GenericParam].IndexSize(), typeDefOrRef})
	tables[ImplMap].SetRowType([6]uint32{2, memberForwarded, stringIndexSize, tables[ModuleRef].IndexSize()})
	tables[InterfaceImpl].SetRowType([6]uint32{tables[TypeDef].IndexSize(), typeDefOrRef})
	tables[ManifestResource].SetRowType([6]uint32{4, 4, stringIndexSize, implementation})
	tables[MemberRef].SetRowType([6]uint32{memberRefParent, stringIndexSize, blobIndexSize})
	tables[MethodDef].SetRowType([6]uint32{4, 2, 2, stringIndexSize, blobIndexSize, tables[Param].IndexSize()})
	tables[MethodImpl].SetRowType([6]uint32{tables[TypeDef].IndexSize(), methodDefOrRef, methodDefOrRef})
	tables[MethodSemantics].SetRowType([6]uint32{2, tables[MethodDef].IndexSize(), hasSemantics})
	tables[MethodSpec].SetRowType([6]uint32{methodDefOrRef, blobIndexSize})
	tables[Module].SetRowType([6]uint32{2, stringIndexSize, guidIndexSize, guidIndexSize, guidIndexSize})
	tables[ModuleRef].SetRowType([6]uint32{stringIndexSize})
	tables[NestedClass].SetRowType([6]uint32{tables[TypeDef].IndexSize(), tables[TypeDef].IndexSize()})
	tables[Param].SetRowType([6]uint32{2, 2, stringIndexSize})
	tables[Property].SetRowType([6]uint32{2, stringIndexSize, blobIndexSize})
	tables[PropertyMap].SetRowType([6]uint32{tables[TypeDef].IndexSize(), tables[Property].IndexSize()})
	tables[StandaloneSig].SetRowType([6]uint32{blobIndexSize})
	tables[TypeDef].SetRowType([6]uint32{
		4,
		stringIndexSize,
		stringIndexSize,
		typeDefOrRef,
		tables[Field].IndexSize(),
		tables[MethodDef].IndexSize(),
	})
	tables[TypeRef].SetRowType([6]uint32{resolutionScope, stringIndexSize, stringIndexSize})
	tables[TypeSpec].SetRowType([6]uint32{blobIndexSize})
}
