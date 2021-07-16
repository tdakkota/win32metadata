package md

import (
	"debug/pe"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	a := require.New(t)

	f, err := pe.Open(`./testdata/.windows/winmd/Windows.Win32.winmd`)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	file, err := ParseMetadata(f)
	a.NoError(err)

	section, err := file.StreamByName("#~")
	a.NoError(err)

	var header TablesHeader
	if err := header.Decode(section); err != nil {
		t.Fatal(err)
	}

	a.Equal([4]byte{}, header.Reserved)
	a.Equal(byte(2), header.MajorVersion)
	a.Equal(byte(0), header.MinorVersion)
}
