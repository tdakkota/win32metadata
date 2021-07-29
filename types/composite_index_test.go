package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/win32metadata/md"
)

func TestCreateIndex(t *testing.T) {
	a := require.New(t)

	idx := CreateHasConstant(md.Param, 10)
	tt, ok := idx.Table()
	a.True(ok)
	a.Equal(md.Param, tt)
	a.Equal(uint32(10), idx.TableIndex())
	a.Panics(func() {
		CreateHasConstant(md.CustomAttribute, 10)
	})
}
