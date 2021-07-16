package md

import (
	"debug/pe"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMetadata(t *testing.T) {
	a := require.New(t)

	f, err := pe.Open(`./testdata/.windows/winmd/Windows.Win32.winmd`)
	a.NoError(err)
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

	f, err := pe.Open(`./testdata/.windows/winmd/Windows.Win32.winmd`)
	a.NoError(err)
	defer f.Close()

	m, err := ParseMetadata(f)
	a.NoError(err)

	_, err = m.StreamByName("#~")
	a.NoError(err)

	_, err = m.StreamByName("lolnogenerics")
	a.Error(err)
}
