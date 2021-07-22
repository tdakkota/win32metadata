package types

import (
	"debug/pe"
	"fmt"
	"io"

	"github.com/tdakkota/win32metadata/md"
)

// Context is a simple helper for accessing file heaps and tables.
type Context struct {
	Metadata *md.Metadata
	section  *io.SectionReader
	md.TablesHeader
}

// FromPE creates new Context from PE file.
func FromPE(f *pe.File) (*Context, error) {
	metadata, err := md.ParseMetadata(f)
	if err != nil {
		return nil, err
	}

	tables, section, err := metadata.Tables()
	if err != nil {
		return nil, err
	}

	return &Context{
		Metadata:     metadata,
		section:      section,
		TablesHeader: tables,
	}, nil
}

// RowCount returns row count of given table type.
func (t *Context) RowCount(tt md.TableType) int {
	return int(t.Tables[tt].RowCount)
}

// Uint32 returns numeric value truncated to uint32.
func (t *Context) Uint32(tt md.TableType, row, column uint32) (uint32, error) {
	return t.Tables[tt].Uint32(t.section, row, column)
}

// Uint64 returns numeric value truncated to uint64.
func (t *Context) Uint64(tt md.TableType, row, column uint32) (uint64, error) {
	return t.Tables[tt].Uint64(t.section, row, column)
}

// String finds string value from #Strings heap using given index column.
func (t *Context) String(tt md.TableType, row, column uint32) (string, error) {
	idx, err := t.Uint32(tt, row, column)
	if err != nil {
		return "", err
	}

	return t.Metadata.ReadString(idx)
}

// Blob finds blob value from #Blob heap using given index column.
func (t *Context) Blob(tt md.TableType, row, column uint32) (Blob, error) {
	idx, err := t.Uint32(tt, row, column)
	if err != nil {
		return nil, err
	}

	return t.Metadata.ReadBlob(idx)
}

// Signature finds signature blob value from #Blob heap using given index column.
func (t *Context) Signature(tt md.TableType, row, column uint32) (Signature, error) {
	sig, err := t.Blob(tt, row, column)
	if err != nil {
		return nil, err
	}

	return Signature(sig), nil
}

// List returns range of indexes using given index.
func (t *Context) List(tt md.TableType, row, column uint32, target md.TableType) (List, error) {
	f, err := t.Uint32(tt, row, column)
	if err != nil {
		return List{}, err
	}
	first := f - 1

	last := t.Tables[target].RowCount
	if row+1 < t.Tables[tt].RowCount {
		l, err := t.Uint32(tt, row+1, column)
		if err != nil {
			return List{}, err
		}
		last = l - 1
	}

	return List{first, last}, nil
}

// Table creates new Table associated with this Context.
func (t *Context) Table(tt md.TableType) Table {
	return Table{
		Type: tt,
		ctx:  t,
	}
}

// Table is a simple helper to access one table.
type Table struct {
	Type md.TableType
	ctx  *Context
}

// RowCount returns row count of this table.
func (t Table) RowCount() int {
	return t.ctx.RowCount(t.Type)
}

// Columns returns columns of this table.
func (t Table) Columns() md.Columns {
	return t.ctx.Tables[t.Type].Columns
}

// Row creates new Row associated with this Table and underlying Context.
func (t Table) Row(row uint32) Row {
	return Row{Table: t, Row: row}
}

// Row is a simple helper to access one table row.
type Row struct {
	Table Table
	Row   uint32
}

// Uint32 returns numeric value truncated to uint32.
func (t *Row) Uint32(column uint32) (uint32, error) {
	return t.Table.ctx.Uint32(t.Table.Type, t.Row, column)
}

// Uint64 returns numeric value truncated to uint64.
func (t *Row) Uint64(column uint32) (uint64, error) {
	return t.Table.ctx.Uint64(t.Table.Type, t.Row, column)
}

// String finds string value from #Strings heap using given index column.
func (t *Row) String(column uint32) (string, error) {
	return t.Table.ctx.String(t.Table.Type, t.Row, column)
}

// Blob finds blob value from #Blob heap using given index column.
func (t *Row) Blob(column uint32) (Blob, error) {
	return t.Table.ctx.Blob(t.Table.Type, t.Row, column)
}

// Signature finds signature blob value from #Blob heap using given index column.
func (t *Row) Signature(column uint32) (Signature, error) {
	return t.Table.ctx.Signature(t.Table.Type, t.Row, column)
}

// List returns range of indexes using given index.
func (t *Row) List(column uint32, target md.TableType) (List, error) {
	return t.Table.ctx.List(t.Table.Type, t.Row, column, target)
}

// ResolveTypeDefOrRefName resolves TypeDefOrRef name.
func (t *Context) ResolveTypeDefOrRefName(ref TypeDefOrRef) (namespace, name string, err error) {
	row, ok := ref.Row(t)
	if !ok || ref.Tag() == 2 /* md.TypeSpec */ {
		return "", "", fmt.Errorf("unexpected tag %v", ref)
	}

	name, err = row.String(1)
	if err != nil {
		return "", "", err
	}

	namespace, err = row.String(2)
	if err != nil {
		return "", "", err
	}

	return namespace, name, nil
}
