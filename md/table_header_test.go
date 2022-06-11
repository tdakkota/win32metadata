package md

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata_Tables(t *testing.T) {
	a := require.New(t)

	f := openTestData(a, "_testdata/Windows.Win32.winmd")
	defer f.Close()

	file, err := ParseMetadata(f)
	a.NoError(err)

	header, section, err := file.Tables()
	a.NoError(err)

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

	{
		fields := []string{
			"value__",
			"ENABLE_ECHO_INPUT",
			"ENABLE_INSERT_MODE",
			"ENABLE_LINE_INPUT",
			"ENABLE_MOUSE_INPUT",
			"ENABLE_PROCESSED_INPUT",
			"ENABLE_QUICK_EDIT_MODE",
			"ENABLE_WINDOW_INPUT",
			"ENABLE_VIRTUAL_TERMINAL_INPUT",
			"ENABLE_PROCESSED_OUTPUT",
		}
		blobs := []struct {
			length    int
			firstByte byte
		}{
			{2, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
			{3, 6},
		}

		field := tables[Field]
		for i, name := range fields {
			idx, err := field.Uint64(section, uint32(i), 1)
			a.NoError(err)
			fieldName, err := file.ReadString(idx)
			a.NoError(err)
			a.Equal(name, fieldName)

			sig, err := field.Uint64(section, uint32(i), 2)
			a.NoError(err)
			blob, err := file.ReadBlob(sig)
			a.NoError(err)
			a.Len(blob, blobs[i].length)
			a.Equal(blob[0], blobs[i].firstByte)
		}
	}
}
