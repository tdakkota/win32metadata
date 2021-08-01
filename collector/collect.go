package collector

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

func (c *collector) Collect(fn func(t Type) error) error {
	typeDefs := c.ctx.Table(md.TypeDef)
	var typeDef types.TypeDef
	for i := uint32(1); i < typeDefs.RowCount(); i++ {
		if err := typeDef.FromRow(typeDefs.Row(i)); err != nil {
			return err
		}

		if typeDef.TypeNamespace == "" ||
			strings.Contains(typeDef.TypeNamespace, "Windows.Win32.System.WinRT") {
			continue
		}

		typ, err := c.collectType(i, typeDef)
		if err != nil {
			return fmt.Errorf(
				"type %s.%s (%d): %w",
				typeDef.TypeNamespace, typeDef.TypeName, i, err,
			)
		}

		if err := fn(typ); err != nil {
			return err
		}
	}

	return nil
}

func (c *collector) collectType(typeIdx types.Index, typeDef types.TypeDef) (Type, error) {
	ns := typeDef.TypeNamespace
	pkg := getPackage(ns)
	typ := Type{
		Name:      publicName(typeDef.TypeName),
		Namespace: ns,
		Package:   pkg,
	}

	imp := &imports{
		currentPkg: pkg,
		types:      map[namedType]Import{},
	}

	fields, err := typeDef.ResolveFieldList(c.ctx)
	if err != nil {
		return typ, err
	}

	{
		fieldsStart := typeDef.FieldList.Start()
		for idx, field := range fields {
			if field.Flags.Static() {
				cst, err := c.collectConstant(imp, fieldsStart+types.Index(idx), field)
				if err != nil {
					return typ, fmt.Errorf("constant %s: %w", field.Name, err)
				}

				typ.Constants = append(typ.Constants, cst)
				continue
			}

			f, err := c.collectField(imp, field)
			if err != nil {
				return typ, fmt.Errorf("field %s: %w", field.Name, err)
			}

			typ.Fields = append(typ.Fields, f)
		}
		typ.IsNewType = len(typ.Fields) == 1
	}

	{
		methodDefs, err := typeDef.ResolveMethodList(c.ctx)
		if err != nil {
			return typ, err
		}

		methodsStart := typeDef.MethodList.Start()
		for idx, methodDef := range methodDefs {
			method, err := c.collectMethod(imp, methodsStart+types.Index(idx), methodDef)
			if err != nil {
				return typ, fmt.Errorf("method %s: %w", methodDef.Name, err)
			}

			typ.Methods = append(typ.Methods, method)
		}
	}

	for _, elem := range imp.types {
		typ.Imports = append(typ.Imports, elem)
	}
	sort.Slice(typ.Imports, func(i, j int) bool {
		return typ.Imports[i].Def < typ.Imports[i].Def &&
			typ.Imports[i].Path < typ.Imports[i].Path
	})

	return typ, nil
}

func (c *collector) collectField(imp *imports, field types.Field) (Field, error) {
	sig, err := field.Signature.Reader().Field(c.ctx)
	if err != nil {
		return Field{}, fmt.Errorf("field signature: %w", err)
	}

	goType, err := c.goType(imp, sig.Field)
	if err != nil {
		return Field{}, err
	}

	return Field{
		Name:   publicName(field.Name),
		GoType: goType,
	}, nil
}

func (c *collector) collectConstant(imp *imports, idx types.Index, field types.Field) (Constant, error) {
	sig, err := field.Signature.Reader().Field(c.ctx)
	if err != nil {
		return Constant{}, fmt.Errorf("const signature: %w", err)
	}

	goType, err := c.goType(imp, sig.Field)
	if err != nil {
		return Constant{}, err
	}

	constantIdx := c.constantIdx[types.CreateHasConstant(md.Field, idx)]
	var cst types.Constant
	if err := cst.FromRow(c.ctx.Table(md.Constant).Row(constantIdx)); err != nil {
		return Constant{}, err
	}

	value, err := readValue(cst)
	if err != nil {
		return Constant{}, fmt.Errorf("value: %w", err)
	}

	return Constant{
		Name:   publicName(field.Name),
		GoType: goType,
		Value:  value,
	}, nil
}

func (c *collector) collectMethod(imp *imports, idx types.Index, method types.MethodDef) (Method, error) {
	sig, err := method.Signature.Reader().Method(c.ctx)
	if err != nil {
		return Method{}, fmt.Errorf("method signature: %w", err)
	}

	params, err := c.resolveParams(method)
	if err != nil {
		return Method{}, fmt.Errorf("resolve params: %w", err)
	}

	returnType, err := c.collectParam(imp, params[0], sig.Return)
	if err != nil {
		return Method{}, fmt.Errorf("return type: %w", err)
	}

	m := Method{
		Name:   publicName(method.Name),
		Return: returnType,
	}
	if implMap, ok := c.implMapIdx[types.CreateMemberForwarded(md.MethodDef, idx)]; ok {
		module, err := implMap.ResolveImportScope(c.ctx)
		if err != nil {
			return Method{}, fmt.Errorf("resolve module: %w", err)
		}

		m.DLLName = module.Name
		m.DLLImport = implMap.ImportName
	}

	for idx, element := range sig.Params {
		p, err := c.collectParam(imp, params[idx+1], element)
		if err != nil {
			return Method{}, fmt.Errorf("param %d: %w", idx, err)
		}

		m.Params = append(m.Params, p)
	}
	return m, nil
}

func (c *collector) resolveParams(method types.MethodDef) (map[int]types.Param, error) {
	paramList, err := method.ResolveParamList(c.ctx)
	if err != nil {
		return nil, err
	}

	params := map[int]types.Param{}
	for _, param := range paramList {
		params[int(param.Sequence)] = param
	}
	return params, nil
}

func (c *collector) collectParam(imp *imports, param types.Param, typ types.Element) (Param, error) {
	goType, err := c.goType(imp, typ)
	if err != nil {
		return Param{}, err
	}

	return Param{
		Name:   cleanupName(param.Name),
		GoType: goType,
	}, nil
}
