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

	tables := header.Tables
	a.Equal(uint32(12), tables[0x0].RowSize)
	a.Equal(uint32(10), tables[0x1].RowSize)
	a.Equal(uint32(22), tables[0x2].RowSize)
	a.Equal(uint32(10), tables[0x4].RowSize)
	a.Equal(uint32(20), tables[0x6].RowSize)
	a.Equal(uint32(8), tables[0x8].RowSize)
	a.Equal(uint32(6), tables[0x9].RowSize)
	a.Equal(uint32(12), tables[0xa].RowSize)
	a.Equal(uint32(10), tables[0xb].RowSize)
	a.Equal(uint32(12), tables[0xc].RowSize)
	a.Equal(uint32(8), tables[0xf].RowSize)
	a.Equal(uint32(8), tables[0x10].RowSize)
	a.Equal(uint32(4), tables[0x1a].RowSize)
	a.Equal(uint32(12), tables[0x1c].RowSize)
	a.Equal(uint32(28), tables[0x20].RowSize)
	a.Equal(uint32(28), tables[0x23].RowSize)
	a.Equal(uint32(4), tables[0x29].RowSize)

	a.Equal(uint32(2), tables[0x0].Columns[0].Size)
	a.Equal(uint32(2), tables[0x1].Columns[0].Size)
	a.Equal(uint32(4), tables[0x2].Columns[0].Size)
	a.Equal(uint32(2), tables[0x4].Columns[0].Size)
	a.Equal(uint32(4), tables[0x6].Columns[0].Size)
	a.Equal(uint32(2), tables[0x8].Columns[0].Size)
	a.Equal(uint32(2), tables[0x9].Columns[0].Size)
	a.Equal(uint32(4), tables[0xa].Columns[0].Size)
	a.Equal(uint32(2), tables[0xb].Columns[0].Size)
	a.Equal(uint32(4), tables[0xc].Columns[0].Size)
	a.Equal(uint32(2), tables[0xf].Columns[0].Size)
	a.Equal(uint32(4), tables[0x10].Columns[0].Size)
	a.Equal(uint32(4), tables[0x1a].Columns[0].Size)
	a.Equal(uint32(2), tables[0x1c].Columns[0].Size)
	a.Equal(uint32(4), tables[0x20].Columns[0].Size)
	a.Equal(uint32(8), tables[0x23].Columns[0].Size)
	a.Equal(uint32(2), tables[0x29].Columns[0].Size)

	v, err := tables[TypeRef].Uint32(section, 0, 0)
	a.NoError(err)
	a.Equal(uint32(4), v)
}
