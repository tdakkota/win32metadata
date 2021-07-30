package collector

import (
	"debug/pe"
	"fmt"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

type typeDefKey struct {
	TypeName      string
	TypeNamespace string
}

type collector struct {
	ctx         *types.Context
	typeNameIdx map[typeDefKey][]types.Index
	// TypeDef index of enclosing class -> map[TypeDef.Name]TypeDef index of nested classes
	nestedIdx map[types.Index]map[string]types.Index
	// MethodDef or Field index -> ImplMap index
	implMapIdx map[types.MemberForwarded]types.ImplMap
	// Param, Field, or Property index -> Constant index
	constantIdx map[types.HasConstant]types.Index
}

func newCollector(f *pe.File) (*collector, error) {
	ctx, err := types.FromPE(f)
	if err != nil {
		return nil, err
	}

	c := &collector{
		ctx: ctx,
	}

	if err := c.readIndex(); err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}
	return c, nil
}

func (c *collector) readIndex() error {
	{
		table := c.ctx.Table(md.TypeDef)
		c.typeNameIdx = make(map[typeDefKey][]types.Index, table.RowCount())

		var typeDef types.TypeDef
		for i := uint32(0); i < table.RowCount(); i++ {
			if err := typeDef.FromRow(table.Row(i)); err != nil {
				return err
			}

			// Skip nested and private types.
			if !typeDef.Flags.Public() {
				continue
			}

			key := typeDefKey{
				TypeName:      typeDef.TypeName,
				TypeNamespace: typeDef.TypeNamespace,
			}
			c.typeNameIdx[key] = append(c.typeNameIdx[key], i)
		}
	}

	{
		table := c.ctx.Table(md.NestedClass)
		c.nestedIdx = make(map[types.Index]map[string]types.Index, table.RowCount())

		var class types.NestedClass
		for i := uint32(0); i < table.RowCount(); i++ {
			if err := class.FromRow(table.Row(i)); err != nil {
				return err
			}
			class.NestedClass--
			class.EnclosingClass--

			nested, err := class.ResolveNestedClass(c.ctx)
			if err != nil {
				return err
			}

			if c.nestedIdx[class.EnclosingClass] == nil {
				c.nestedIdx[class.EnclosingClass] = map[string]types.Index{}
			}
			c.nestedIdx[class.EnclosingClass][nested.TypeName] = class.NestedClass
		}
	}

	{
		table := c.ctx.Table(md.ImplMap)
		c.implMapIdx = make(map[types.MemberForwarded]types.ImplMap, table.RowCount())

		var implMap types.ImplMap
		for i := uint32(0); i < table.RowCount(); i++ {
			if err := implMap.FromRow(table.Row(i)); err != nil {
				return err
			}

			c.implMapIdx[implMap.MemberForwarded] = implMap
		}
	}

	{
		table := c.ctx.Table(md.Constant)
		c.constantIdx = make(map[types.HasConstant]types.Index, table.RowCount())

		var constant types.Constant
		for i := uint32(0); i < table.RowCount(); i++ {
			if err := constant.FromRow(table.Row(i)); err != nil {
				return err
			}

			c.constantIdx[constant.Parent] = i
		}
	}
	return nil
}
