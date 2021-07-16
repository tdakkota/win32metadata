package md

import (
	"bytes"
	"io"
	"strings"
)

// MetadataRoot is a representation of II.24.2.1 Metadata root
type MetadataRoot struct {
	MajorVersion  uint16
	MinorVersion  uint16
	Reserved      [4]byte
	Version       string
	Flags         uint16
	StreamHeaders []StreamHeader
}

func (m *MetadataRoot) Decode(r io.Reader) error {
	rr := reader{r: r}

	// II.24.2.1 Metadata root
	//
	// Field Length.
	// Number of bytes allocated to hold version string (including null terminator), call this x.
	// Call the length of the string (including the terminator) m (we require m <= 255);
	// the length x is m rounded up to a multiple of four.
	var versionLength uint32

	if !rr.Read(&m.MajorVersion) ||
		!rr.Read(&m.MinorVersion) ||
		!rr.Read(&m.Reserved) ||
		!rr.Read(&versionLength) {
		return rr.Err()
	}

	{
		b := &strings.Builder{}
		b.Grow(int(versionLength))

		if _, err := io.CopyN(b, r, int64(versionLength)); err != nil {
			return err
		}
		m.Version = strings.TrimRight(b.String(), "\x00")
	}

	// Number of streams, say n.
	var streamsN uint16
	if !rr.Read(&m.Flags) || !rr.Read(&streamsN) {
		return rr.Err()
	}

	var hdr StreamHeader
	for i := 0; i < int(streamsN); i++ {
		if err := hdr.Decode(r); err != nil {
			return err
		}
		m.StreamHeaders = append(m.StreamHeaders, hdr)
	}

	return nil
}

// StreamHeader is a representation of II.24.2.2 Stream header.
type StreamHeader struct {
	Offset uint32
	Size   uint32
	Name   string
}

func (h *StreamHeader) Decode(r io.Reader) error {
	rr := reader{r: r}

	if !rr.Read(&h.Offset) || !rr.Read(&h.Size) {
		return rr.Err()
	}

	// II.24.2.2 Stream header
	//
	// Field "Name".
	// Name of the stream as null-terminated variable length array
	// of ASCII characters, padded to the next 4-byte boundary
	// with \0 characters. The name is limited to 32 characters
	var (
		b   = &strings.Builder{}
		buf = make([]byte, 4)
	)
	for i := 0; i < 32; i++ {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return err
		}

		idx := bytes.IndexByte(buf, 0)
		if idx >= 0 {
			b.Write(buf[:idx])
			break
		}
		b.Write(buf)
	}
	h.Name = b.String()

	return nil
}
