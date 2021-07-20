package md

import (
	"debug/pe"
	"fmt"
	"io"
	"strings"
)

// Metadata is a simple wrapper around MetadataRoot to access
// metadata streams.
type Metadata struct {
	r *io.SectionReader
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
	heap, ok := m.findStreamHeader("#Strings")
	if !ok {
		return "", fmt.Errorf("string heap section not found")
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

	return buf.String(), nil
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
		MetadataRoot: root,
	}, nil
}
