package collector

import (
	"fmt"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

func (c *collector) resolveTypeRef(row types.Row) ([]types.Index, error) {
	var typeRef types.TypeRef
	if err := typeRef.FromRow(row); err != nil {
		return nil, err
	}

	scope := typeRef.ResolutionScope
	if tt, ok := scope.Table(); ok && tt == md.TypeRef {
		sub, ok := scope.Row(c.ctx)
		if !ok {
			return nil, fmt.Errorf("unexpected tag %v", scope)
		}

		typeDefs, err := c.resolveTypeRef(sub)
		if err != nil {
			return nil, err
		}

		for _, idx := range typeDefs {
			if classes, ok := c.nestedIdx[idx]; ok {
				if typeDefIdx, ok := classes[typeRef.TypeName]; ok {
					return []types.Index{typeDefIdx}, err
				}
			}
		}

		return nil, fmt.Errorf(
			"type ref %s.%s (scope: %v) not found",
			typeRef.TypeNamespace, typeRef.TypeName, typeRef.ResolutionScope,
		)
	}

	if v, ok := c.typeNameIdx[typeDefKey{
		TypeName:      typeRef.TypeName,
		TypeNamespace: typeRef.TypeNamespace,
	}]; ok {
		return v, nil
	}

	return nil, fmt.Errorf(
		"type ref %s.%s %v not found",
		typeRef.TypeNamespace, typeRef.TypeName, typeRef.ResolutionScope,
	)
}

func (c *collector) resolveTypeDef(ref types.TypeDefOrRef) (types.TypeDef, error) {
	tt, ok := ref.Table()
	if !ok {
		return types.TypeDef{}, fmt.Errorf("unexpected tag %v", ref)
	}

	switch tt {
	case md.TypeDef:
		return c.typeDef(ref.TableIndex())
	case md.TypeRef:
		idx, err := c.resolveTypeRef(c.ctx.Table(md.TypeRef).Row(ref.TableIndex()))
		if err != nil {
			return types.TypeDef{}, err
		}
		if len(idx) < 1 {
			return types.TypeDef{}, fmt.Errorf("can't resolve %v", ref)
		}

		return c.typeDef(idx[0])
	default:
		return types.TypeDef{}, fmt.Errorf("unexpected table type: %d", tt)
	}
}

func (c *collector) typeDef(idx types.Index) (types.TypeDef, error) {
	var v types.TypeDef
	if err := v.FromRow(c.ctx.Table(md.TypeDef).Row(idx)); err != nil {
		return v, err
	}
	return v, nil
}
