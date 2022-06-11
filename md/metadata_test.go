package md

import (
	"bytes"
	"debug/pe"
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed _testdata
	testdata embed.FS
)

func openTestData(a *require.Assertions, p string) *pe.File {
	data, err := testdata.ReadFile(p)
	a.NoError(err)

	f, err := pe.NewFile(bytes.NewReader(data))
	a.NoError(err)

	return f
}

func TestParseMetadata(t *testing.T) {
	a := require.New(t)

	f := openTestData(a, `_testdata/Windows.Win32.winmd`)
	defer f.Close()

	m, err := ParseMetadata(f)
	a.NoError(err)

	// II.24.2.1 Metadata root
	//
	// Major version, 1 (ignore on read)
	a.Equal(uint16(1), m.MajorVersion)
	// Minor version, 1 (ignore on read)
	a.Equal(uint16(1), m.MinorVersion)
	// Reserved, always 0 (§II.24.1).
	a.Equal([4]byte{}, m.Reserved)
	// Reserved, always 0 (§II.24.1).
	a.Equal(uint16(0), m.Flags)
	a.NotEmpty(m.StreamHeaders)
}

func TestMetadata_StreamByName(t *testing.T) {
	a := require.New(t)

	f := openTestData(a, `_testdata/Windows.WinRT.winmd`)
	defer f.Close()

	m, err := ParseMetadata(f)
	a.NoError(err)

	_, err = m.StreamByName("#~")
	a.NoError(err)

	_, err = m.StreamByName("lolnogenerics")
	a.Error(err)
}
