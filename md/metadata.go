package md

import (
	"debug/pe"
	"fmt"
	"io"
)

// Metadata is a simple wrapper around MetadataRoot to access
// metadata streams.
type Metadata struct {
	r *io.SectionReader
	MetadataRoot
}

// StreamByName finds and returns metadata section reader by name.
func (m Metadata) StreamByName(name string) (*io.SectionReader, error) {
	var section *io.SectionReader
	for _, shr := range m.StreamHeaders {
		if shr.Name == name {
			section = io.NewSectionReader(m.r, int64(shr.Offset), int64(shr.Size))
		}
	}
	if section == nil {
		return nil, fmt.Errorf("section %q not found", name)
	}

	return section, nil
}

// ParseMetadata parses and creates Metadata from given PE file.
func ParseMetadata(f *pe.File) (Metadata, error) {
	cliHeader, err := getCLIHeader(f)
	if err != nil {
		return Metadata{}, err
	}

	r, err := getMetadataReader(f, cliHeader)
	if err != nil {
		return Metadata{}, err
	}

	var root MetadataRoot
	if err := root.Decode(r); err != nil {
		return Metadata{}, err
	}

	return Metadata{
		r:            r,
		MetadataRoot: root,
	}, nil
}
