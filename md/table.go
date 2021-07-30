package md

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TableType is a metadata table type.
//
// See II.22.1 Metadata validation rules.
type TableType int

//go:generate go run golang.org/x/tools/cmd/stringer -type=TableType

const (
	// Module table type.
	Module TableType = 0x00
	// TypeRef table type.
	TypeRef TableType = 0x01
	// TypeDef table type.
	TypeDef TableType = 0x02
	// Field table type.
	Field TableType = 0x04
	// MethodDef table type.
	MethodDef TableType = 0x06
	// Param table type.
	Param TableType = 0x08
	// InterfaceImpl table type.
	InterfaceImpl TableType = 0x09
	// MemberRef table type.
	MemberRef TableType = 0x0a
	// Constant table type.
	Constant TableType = 0x0b
	// CustomAttribute table type.
	CustomAttribute TableType = 0x0c
	// FieldMarshal table type.
	FieldMarshal TableType = 0x0d
	// DeclSecurity table type.
	DeclSecurity TableType = 0x0e
	// ClassLayout table type.
	ClassLayout TableType = 0x0f
	// FieldLayout table type.
	FieldLayout TableType = 0x10
	// StandAloneSig table type.
	StandAloneSig TableType = 0x11
	// EventMap table type.
	EventMap TableType = 0x12
	// Event table type.
	Event TableType = 0x14
	// PropertyMap table type.
	PropertyMap TableType = 0x15
	// Property table type.
	Property TableType = 0x17
	// MethodSemantics table type.
	MethodSemantics TableType = 0x18
	// MethodImpl table type.
	MethodImpl TableType = 0x19
	// ModuleRef table type.
	ModuleRef TableType = 0x1a
	// TypeSpec table type.
	TypeSpec TableType = 0x1b
	// ImplMap table type.
	ImplMap TableType = 0x1c
	// FieldRva table type.
	FieldRva TableType = 0x1d
	// Assembly table type.
	Assembly TableType = 0x20
	// AssemblyProcessor table type.
	AssemblyProcessor TableType = 0x21
	// AssemblyOs table type.
	AssemblyOs TableType = 0x22
	// AssemblyRef table type.
	AssemblyRef TableType = 0x23
	// AssemblyRefProcessor table type.
	AssemblyRefProcessor TableType = 0x24
	// AssemblyRefOs table type.
	AssemblyRefOs TableType = 0x25
	// File table type.
	File TableType = 0x26
	// ExportedType table type.
	ExportedType TableType = 0x27
	// ManifestResource table type.
	ManifestResource TableType = 0x28
	// NestedClass table type.
	NestedClass TableType = 0x29
	// GenericParam table type.
	GenericParam TableType = 0x2a
	// MethodSpec table type.
	MethodSpec TableType = 0x2b
	// GenericParamConstraint table type.
	GenericParamConstraint TableType = 0x2c
)

// Column represents metadata table column sizes.
type Column struct {
	// Offset from row start.
	Offset uint32
	// Size of column in bytes.
	Size uint32
}

// Zero denotes that column has zero value.
func (c Column) Zero() bool {
	return c == Column{}
}

// Columns is a table columns list.
type Columns = [6]Column

// Table represents metadata table header.
type Table struct {
	Type     TableType
	Offset   int64
	RowCount uint32
	RowSize  uint32
	Columns  Columns
}

// Find returns offset of given column in row.
func (t Table) Find(row, column uint32) (uint32, error) {
	c := t.Columns[column]
	if c.Zero() {
		return 0, fmt.Errorf("type %#x does not have column %d", t.Type, column)
	}
	if row > t.RowCount {
		return 0, fmt.Errorf("row index %d is out of bounds (%d)", row, t.RowCount)
	}

	return uint32(t.Offset) + row*t.RowSize + c.Offset, nil
}

// Uint32 returns numeric value truncated to uint32.
func (t Table) Uint32(r io.ReaderAt, row, column uint32) (uint32, error) {
	offset, err := t.Find(row, column)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, t.Columns[column].Size)
	if _, err := r.ReadAt(buf, int64(offset)); err != nil {
		return 0, err
	}

	switch len(buf) {
	case 1:
		return uint32(buf[0]), nil
	case 2:
		return uint32(binary.LittleEndian.Uint16(buf)), nil
	case 4:
		return binary.LittleEndian.Uint32(buf), nil
	default:
		return uint32(binary.LittleEndian.Uint64(buf)), nil
	}
}

// Uint64 returns numeric value truncated to uint64.
func (t Table) Uint64(r io.ReaderAt, row, column uint32) (uint64, error) {
	offset, err := t.Find(row, column)
	if err != nil {
		return 0, err
	}
	buf := make([]byte, t.Columns[column].Size)
	if _, err := r.ReadAt(buf, int64(offset)); err != nil {
		return 0, err
	}

	switch len(buf) {
	case 1:
		return uint64(buf[0]), nil
	case 2:
		return uint64(binary.LittleEndian.Uint16(buf)), nil
	case 4:
		return uint64(binary.LittleEndian.Uint32(buf)), nil
	default:
		return binary.LittleEndian.Uint64(buf), nil
	}
}

// IndexSize returns size of table index.
func (t Table) IndexSize() uint32 {
	if t.RowCount < (1 << 16) {
		return 2
	}

	return 4
}

// SetRowType sets rows sizes.
// NB: Zero size means that column is not present.
func (t *Table) SetRowType(sizes [6]uint32) {
	t.RowSize = 0
	for i, column := range sizes {
		if column == 0 {
			break
		}
		t.Columns[i] = Column{
			Offset: t.RowSize,
			Size:   column,
		}
		t.RowSize += column
	}
}
