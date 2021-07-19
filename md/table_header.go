package md

import (
	"io"
)

type RowCount struct {
	Tag   TableType
	Count uint32
}

// TablesHeader is a representation of II.24.2.6 #~ stream.
type TablesHeader struct {
	Reserved     [4]byte
	MajorVersion uint8
	MinorVersion uint8
	HeapSizes    uint8
	Reserved2    [1]byte
	Valid        uint64
	Sorted       uint64
	Tables       [GenericParamConstraint + 1]Table
}

// StringIndexSize returns size of index of "#String" heap.
func (h TablesHeader) StringIndexSize() uint32 {
	if h.HeapSizes&1 == 1 {
		return 4
	}
	return 2
}

// BlobIndexSize returns size of index of "#Blob" heap.
func (h TablesHeader) BlobIndexSize() uint32 {
	if (h.HeapSizes>>2)&1 == 1 {
		return 4
	}
	return 2
}

// GUIDIndexSize returns size of index of "#GUID" heap.
func (h TablesHeader) GUIDIndexSize() uint32 {
	if (h.HeapSizes>>1)&1 == 1 {
		return 4
	}
	return 2
}

// Decode decodes TablesHeader from stream.
func (h *TablesHeader) Decode(r io.Reader) error {
	rr := reader{r: r}

	if !rr.Read(&h.Reserved) ||
		!rr.Read(&h.MajorVersion) ||
		!rr.Read(&h.MinorVersion) ||
		!rr.Read(&h.HeapSizes) ||
		!rr.Read(&h.Reserved2) ||
		!rr.Read(&h.Valid) ||
		!rr.Read(&h.Sorted) {
		return rr.Err()
	}

	// II.24.2.6 #~ stream
	// ...
	// The Valid field is a 64-bit bitvector that has a specific bit set for each table that is stored in the stream;
	// the mapping of tables to indexes is given at the start of §II.22.
	var (
		row    uint32
		offset int64 = 24 // Constant structure size
	)
	for i := 0; i < 64; i++ {
		if h.Valid>>i&1 == 0 {
			continue
		}

		if !rr.Read(&row) {
			return rr.Err()
		}
		offset += 4
		h.Tables[i] = Table{
			Type:     TableType(i),
			RowCount: row,
		}
	}
	h.computeIndexes()

	// Compute data offsets of every table.
	{
		for i, table := range h.Tables {
			h.Tables[i].Offset = offset
			offset += int64(table.RowCount * table.RowSize)
		}
	}
	return nil
}
