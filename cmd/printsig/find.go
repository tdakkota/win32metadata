package main

import (
	"fmt"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

func findMethod(c *types.Context, typeNamespace, methodName string) (uint32, types.MethodDef, error) {
	checkNamespace := typeNamespace != ""

	typeDefs := c.Table(md.TypeDef)
	methodDefs := c.Table(md.MethodDef)

	var typeDef types.TypeDef
	for i := uint32(0); i < typeDefs.RowCount(); i++ {
		if err := typeDef.FromRow(typeDefs.Row(i)); err != nil {
			return 0, types.MethodDef{}, err
		}

		if checkNamespace && typeDef.TypeNamespace != typeNamespace {
			continue
		}

		list := typeDef.MethodList
		if list.Empty() {
			continue
		}

		var methodDef types.MethodDef
		for methodIdx := list.Start(); methodIdx < list.End(); methodIdx++ {
			if err := methodDef.FromRow(methodDefs.Row(methodIdx)); err != nil {
				return 0, types.MethodDef{}, err
			}

			if methodName == methodDef.Name {
				return methodIdx, methodDef, nil
			}
		}
	}

	return 0, types.MethodDef{}, fmt.Errorf("method %q not found", methodName)
}

func resolveTypeRef(t *types.Context, row types.Row) (types.TypeDef, uint32, error) {
	var r types.TypeRef
	if err := r.FromRow(row); err != nil {
		return types.TypeDef{}, 0, err
	}

	scope := r.ResolutionScope
	if tt, ok := scope.Table(); ok && tt == md.TypeRef {
		sub, ok := scope.Row(t)
		if !ok {
			return types.TypeDef{}, 0, fmt.Errorf("unexpected tag %v", scope)
		}

		_, idx, err := resolveTypeRef(t, sub)
		if err != nil {
			return types.TypeDef{}, 0, err
		}

		nestedClasses := t.Table(md.NestedClass)
		var class types.NestedClass
		for i := uint32(0); i < nestedClasses.RowCount(); i++ {
			if err := class.FromRow(nestedClasses.Row(i)); err != nil {
				return types.TypeDef{}, 0, err
			}

			if class.EnclosingClass == idx {
				def, err := class.ResolveNestedClass(t)
				if err != nil {
					return types.TypeDef{}, 0, err
				}

				return def, class.NestedClass, nil
			}
		}
	}

	typeDefs := t.Table(md.TypeDef)
	var def types.TypeDef
	for i := uint32(0); i < typeDefs.RowCount(); i++ {
		if err := def.FromRow(typeDefs.Row(i)); err != nil {
			return types.TypeDef{}, 0, err
		}
		if r.TypeName == def.TypeName && r.TypeNamespace == def.TypeNamespace {
			return def, i, nil
		}
	}

	return types.TypeDef{}, 0, fmt.Errorf("TypeDef %s %s not found", r.TypeNamespace, r.TypeName)
}

func resolveTypeDef(t *types.Context, ref types.TypeDefOrRef) (types.TypeDef, error) {
	// TODO(tdakkota): cache find TypeDef or build an index.

	row, ok := ref.Row(t)
	if !ok {
		return types.TypeDef{}, fmt.Errorf("unexpected tag %v", ref)
	}

	switch tt := row.Table.Type; tt {
	case md.TypeDef:
		var def types.TypeDef
		if err := def.FromRow(row); err != nil {
			return def, err
		}
		return def, nil
	case md.TypeRef:
		def, _, err := resolveTypeRef(t, row)
		if err != nil {
			return types.TypeDef{}, err
		}
		return def, nil
	default:
		return types.TypeDef{}, fmt.Errorf("unexpected table type: %d", tt)
	}
}
