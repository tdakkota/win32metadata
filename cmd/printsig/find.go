package main

import (
	"fmt"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

func findMethod(c *types.Context, typeNamespace, methodName string) (types.MethodDef, error) {
	table := c.Table(md.TypeDef)
	var typeDef types.TypeDef

	checkNamespace := typeNamespace != ""
	for i := uint32(0); i < table.RowCount(); i++ {
		if err := typeDef.FromRow(table.Row(i)); err != nil {
			return types.MethodDef{}, err
		}

		if checkNamespace && typeDef.TypeNamespace != typeNamespace {
			continue
		}

		if typeDef.MethodList.Empty() {
			continue
		}

		methods, err := typeDef.ResolveMethodList(c)
		if err != nil {
			return types.MethodDef{}, err
		}

		for _, method := range methods {
			if methodName != method.Name {
				continue
			}

			return method, nil
		}
	}

	return types.MethodDef{}, fmt.Errorf("method %q not found", methodName)
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
			class.NestedClass--
			class.EnclosingClass--

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
