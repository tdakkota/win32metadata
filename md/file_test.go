package md

import (
	"debug/pe"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	a := require.New(t)

	f, err := pe.Open(`./testdata/.windows/winmd/Windows.Win32.winmd`)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	header, err := getCLIHeader(f)
	if err != nil {
		t.Fatal(err)
	}
	a.Equal(int(unsafe.Sizeof(header)), int(header.CB))

	r, err := getMetadataReader(f, header)
	if err != nil {
		t.Fatal(err)
	}

	var file File
	if err := file.Decode(r); err != nil {
		t.Fatal(err)
	}

	a.Equal(uint16(1), file.MajorVersion)
	a.Equal(uint16(1), file.MinorVersion)
	a.Equal("v4.0.30319", file.Version)
}
