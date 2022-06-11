package win32metadata

import (
	"bytes"
	"debug/pe"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

var (
	//go:embed md/_testdata/Windows.Win32.winmd
	win32 []byte
)

func TestIssue33(t *testing.T) {
	a := require.New(t)

	f, err := pe.NewFile(bytes.NewReader(win32))
	a.NoError(err)
	defer f.Close()

	c, err := types.FromPE(f)
	a.NoError(err)

	tt := c.Table(md.InterfaceImpl)
	var row types.InterfaceImpl
	for i := uint32(0); i < tt.RowCount(); i++ {
		a.NoError(row.FromRow(tt.Row(i)))

		class, err := row.ResolveClass(c)
		a.NoError(err)

		// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitemarray
		if class.TypeNamespace != "Windows.Win32.UI.Shell" ||
			class.TypeName != "IShellItemArray" {
			continue
		}

		ns, name, err := c.ResolveTypeDefOrRefName(row.Interface)
		a.NoError(err)
		a.Equal("Windows.Win32.System.Com", ns)
		a.Equal("IUnknown", name)

		methods, err := class.ResolveMethodList(c)
		a.NoError(err)

		var names []string
		for _, m := range methods {
			names = append(names, m.Name)
		}
		a.Equal([]string{
			"BindToHandler",
			"GetPropertyStore",
			"GetPropertyDescriptionList",
			"GetAttributes",
			"GetCount",
			"GetItemAt",
			"EnumItems",
		}, names)

		return
	}
	t.Fatal("Can't find Windows.Win32.UI.Shell.IShellItemArray")
}
