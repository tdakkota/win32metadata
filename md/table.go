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

const (
	Module                 TableType = 0x00
	TypeRef                TableType = 0x01
	TypeDef                TableType = 0x02
	Field                  TableType = 0x04
	MethodDef              TableType = 0x06
	Param                  TableType = 0x08
	InterfaceImpl          TableType = 0x09
	MemberRef              TableType = 0x0a
	Constant               TableType = 0x0b
	CustomAttribute        TableType = 0x0c
	FieldMarshal           TableType = 0x0d
	DeclSecurity           TableType = 0x0e
	ClassLayout            TableType = 0x0f
	FieldLayout            TableType = 0x10
	StandaloneSig          TableType = 0x11
	EventMap               TableType = 0x12
	Event                  TableType = 0x14
	PropertyMap            TableType = 0x15
	Property               TableType = 0x17
	MethodSemantics        TableType = 0x18
	MethodImpl             TableType = 0x19
	ModuleRef              TableType = 0x1a
	TypeSpec               TableType = 0x1b
	ImplMap                TableType = 0x1c
	FieldRva               TableType = 0x1d
	Assembly               TableType = 0x20
	AssemblyProcessor      TableType = 0x21
	AssemblyOs             TableType = 0x22
	AssemblyRef            TableType = 0x23
	AssemblyRefProcessor   TableType = 0x24
	AssemblyRefOs          TableType = 0x25
	File                   TableType = 0x26
	ExportedType           TableType = 0x27
	ManifestResource       TableType = 0x28
	NestedClass            TableType = 0x29
	GenericParam           TableType = 0x2a
	MethodSpec             TableType = 0x2b
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

// Table represents metadata table header.
type Table struct {
	Type     TableType
	Offset   int64
	RowCount uint32
	RowSize  uint32
	Columns  [6]Column
}

// Find returns offset of given column in row.
func (t Table) Find(row, column uint32) uint32 {
	c := t.Columns[column]
	if c.Zero() {
		panic(fmt.Sprintf("type %#x does not have column %d", t.Type, column))
	}
	if row >= t.RowCount {
		panic(fmt.Sprintf("row index %d is out of bounds (%d)", row, t.RowCount))
	}

	return uint32(t.Offset) + row*t.RowSize + c.Offset
}

// Uint32 returns numeric value truncated to uint32.
func (t Table) Uint32(r io.ReaderAt, row, column uint32) (uint32, error) {
	offset := t.Find(row, column)
	buf := make([]byte, t.Columns[column].Size, 8)
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

// IndexSize returns size of table index.
func (t Table) IndexSize() uint32 {
	if t.RowCount < (1 << 16) {
		return 2
	} else {
		return 4
	}
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