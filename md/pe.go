package md

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

// CLIHeader is a II.25.3.3 CLI header representation.
type CLIHeader struct {
	CB                             uint32
	MajorRuntimeVersion            uint16
	MinorRuntimeVersion            uint16
	MetaData                       pe.DataDirectory
	Flags                          uint32
	EntryPointTokenOrEntryPointRva uint32
	Resources                      pe.DataDirectory
	StrongNameSignature            pe.DataDirectory
	CodeManagerTable               pe.DataDirectory
	VtableFixups                   pe.DataDirectory
	ExportAddressTableJumps        pe.DataDirectory
	ManagedNativeHeader            pe.DataDirectory
}

func findSection(sections []*pe.Section, va uint32) (*pe.Section, error) {
	var section *pe.Section
	for _, s := range sections {
		if va >= s.VirtualAddress && va < s.VirtualAddress+s.VirtualSize {
			section = s
			break
		}
	}
	if section == nil {
		return nil, fmt.Errorf("metadata section not found (%d sections checked)", len(sections))
	}

	return section, nil
}

func getCLIHeader(f *pe.File) (CLIHeader, error) {
	// See II.25.2.3.3 PE header data directories.
	const IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR = (208 - 96) / 8 // 14

	var (
		comVirtualAddress uint32
	)
	switch v := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		comVirtualAddress = v.DataDirectory[IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR].VirtualAddress
	case *pe.OptionalHeader64:
		comVirtualAddress = v.DataDirectory[IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR].VirtualAddress
	default:
		return CLIHeader{}, fmt.Errorf("unexpected type %T", v)
	}

	section, err := findSection(f.Sections, comVirtualAddress)
	if err != nil {
		return CLIHeader{}, err
	}

	const HeaderSize = int64(unsafe.Sizeof(CLIHeader{}))
	headerReader := io.NewSectionReader(
		section,
		int64(comVirtualAddress-section.VirtualAddress),
		HeaderSize, // Size of CLI header
	)

	var h CLIHeader
	if err := binary.Read(headerReader, binary.LittleEndian, &h); err != nil {
		return CLIHeader{}, err
	}
	if int64(h.CB) != HeaderSize {
		return CLIHeader{}, fmt.Errorf("invalid size of CLI header: %d", h.CB)
	}

	return h, err
}

func getMetadataReader(f *pe.File, header CLIHeader) (*io.SectionReader, error) {
	section, err := findSection(f.Sections, header.MetaData.VirtualAddress)
	if err != nil {
		return nil, err
	}
	r := io.NewSectionReader(
		section,
		int64(header.MetaData.VirtualAddress-section.VirtualAddress),
		int64(header.MetaData.Size),
	)

	var (
		magic uint32
	)
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return nil, err
	}

	const STORAGE_MAGIC_SIG = 0x424A_5342
	if magic != STORAGE_MAGIC_SIG {
		return nil, fmt.Errorf("invalid magic: %x, expected %x", magic, STORAGE_MAGIC_SIG)
	}

	return r, nil
}
