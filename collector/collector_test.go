package collector

import (
	"bytes"
	"debug/pe"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

type Writer struct {
	dir string
}

func (w Writer) Write(namespace, typeName string, data []byte) error {
	namespaces := getNamespaces(namespace)

	path := filepath.Join(namespaces...)
	path = filepath.Join(w.dir, path)
	if err := os.MkdirAll(path, 0775); err != nil {
		return err
	}
	path = filepath.Join(path, typeName+".go")

	formatted, err := format.Source(data)
	if err != nil {
		os.WriteFile(path+".dump", data, 0775)
		return fmt.Errorf("format %q: %w", path, err)
	}

	return os.WriteFile(path, formatted, 0775)
}

func TestGenerateAll(t *testing.T) {
	a := require.New(t)

	tmpl, err := template.ParseFiles("gen.tmpl")
	if err != nil {
		a.NoError(err)
	}
	w := Writer{dir: "mdgen"}

	file, err := pe.Open(`../md/testdata/.windows/winmd/Windows.Win32.winmd`)
	if err != nil {
		a.NoError(err)
	}
	defer func() {
		_ = file.Close()
	}()

	col, err := newCollector(file)
	if err != nil {
		a.NoError(err)
	}

	buf := bytes.Buffer{}
	a.NoError(col.Collect(func(t Type) error {
		if t.Name == "<Module>" {
			return nil
		}

		buf.Reset()
		if err := tmpl.Execute(&buf, t); err != nil {
			return fmt.Errorf("execute: %w", err)
		}

		return w.Write(t.Namespace, t.Name, buf.Bytes())
	}))
}
