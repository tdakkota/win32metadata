package md

import (
	"debug/pe"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := pe.Open(`./.windows/winmd/Windows.Win32.winmd`)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	header, err := getCLIHeader(f)
	if err != nil {
		t.Fatal(err)
	}

	r, err := getMetadataReader(f, header)
	if err != nil {
		t.Fatal(err)
	}

	var file File
	if err := file.Decode(r); err != nil {
		t.Fatal(err)
	}

	for _, hdr := range file.StreamHeaders {
		fmt.Println(hdr.Name)
	}
}
