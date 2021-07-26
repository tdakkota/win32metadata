package md

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Metadata is a simple wrapper around MetadataRoot to access
// metadata streams.
type Metadata struct {
	r *io.SectionReader
	// Cache strings from file to prevent allocations.
	strings map[uint32]string
	MetadataRoot
}

func (m *Metadata) findStreamHeader(name string) (StreamHeader, bool) {
	for _, shr := range m.StreamHeaders {
		if shr.Name == name {
			return shr, true
		}
	}
	return StreamHeader{}, false
}

// StreamByName finds and returns metadata section reader by name.
func (m *Metadata) StreamByName(name string) (*io.SectionReader, error) {
	shr, ok := m.findStreamHeader(name)
	if !ok {
		return nil, fmt.Errorf("section %q not found", name)
	}

	return io.NewSectionReader(m.r, int64(shr.Offset), int64(shr.Size)), nil
}

// Tables decodes metadata tables header and returns it and table data reader.
func (m *Metadata) Tables() (TablesHeader, *io.SectionReader, error) {
	section, err := m.StreamByName("#~")
	if err != nil {
		return TablesHeader{}, nil, err
	}

	var header TablesHeader
	if err := header.Decode(section); err != nil {
		return TablesHeader{}, nil, err
	}

	return header, section, nil
}

// ReadString reads string from String heap.
func (m *Metadata) ReadString(idx uint32) (string, error) {
	if v, ok := m.strings[idx]; ok {
		return v, nil
	}

	heap, ok := m.findStreamHeader("#Strings")
	if !ok {
		return "", fmt.Errorf("string heap stream not found")
	}

	var (
		offset = int64(heap.Offset + idx)
		one    [1]byte
		buf    strings.Builder
	)
	for {
		_, err := m.r.ReadAt(one[:], offset)
		if err != nil {
			return "", err
		}
		if one[0] == '\x00' {
			break
		}

		buf.WriteByte(one[0])
		offset++
	}

	v := buf.String()
	m.strings[idx] = v
	return v, nil
}

// ReadBlob reads blob from Blob heap.
func (m *Metadata) ReadBlob(idx uint32) ([]byte, error) {
	// TODO(tdakkota): Decode blob lazily using io.Reader/some helper.

	heap, ok := m.findStreamHeader("#Blob")
	if !ok {
		return nil, fmt.Errorf("blob heap stream not found")
	}
	var (
		offset = int64(heap.Offset + idx)
		buf    = make([]byte, 4)
		// Size of blob data
		blobSize int
		// Size of length in bytes
		lenSize int64
	)

	_, err := m.r.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}

	switch v := buf[0] >> 5; {
	case v <= 3:
		lenSize = 1
		blobSize = int(buf[0] & 0x7f)
	case v >= 4 && v <= 5:
		lenSize = 2
		blobSize = int(binary.LittleEndian.Uint16([]byte{buf[0] & 0x3f, buf[1]}))
	case v == 6:
		lenSize = 4
		blobSize = int(binary.LittleEndian.Uint32([]byte{buf[0] & 0x1f, buf[1], buf[2], buf[3]}))
	default:
		return nil, fmt.Errorf("invalid blob length: %d", buf[0])
	}

	buf = append(buf[:0], make([]byte, blobSize)...)
	if _, err := m.r.ReadAt(buf, offset+lenSize); err != nil {
		return nil, err
	}

	return buf, nil
}

// ParseMetadata parses and creates Metadata from given PE file.
func ParseMetadata(f *pe.File) (*Metadata, error) {
	cliHeader, err := getCLIHeader(f)
	if err != nil {
		return nil, err
	}

	r, err := getMetadataReader(f, cliHeader)
	if err != nil {
		return nil, err
	}

	var root MetadataRoot
	if err := root.Decode(r); err != nil {
		return nil, err
	}

	return &Metadata{
		r:            r,
		strings:      map[uint32]string{},
		MetadataRoot: root,
	}, nil
}
