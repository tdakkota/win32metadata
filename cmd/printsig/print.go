package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tdakkota/win32metadata/types"
)

func queueTypeDefs(
	c *types.Context,
	e types.Element,
	toPrint map[types.TypeDefOrRef]types.TypeDef,
) (namespace, name string, err error) {
	idx := e.Type.TypeDef.Index

	def, ok := toPrint[idx]
	if ok {
		return def.TypeNamespace, def.TypeName, nil
	}

	d, err := resolveTypeDef(c, idx)
	if err != nil {
		return "", "", err
	}

	fieldList, err := d.ResolveFieldList(c)
	if err != nil {
		return "", "", err
	}

	for _, field := range fieldList {
		sig, err := field.Signature.Reader().Field(c)
		if err != nil {
			return "", "", err
		}

		kind := sig.Field.Type.Kind
		if kind != types.ELEMENT_TYPE_VALUETYPE &&
			kind != types.ELEMENT_TYPE_CLASS {
			continue
		}

		if _, _, err := queueTypeDefs(c, sig.Field, toPrint); err != nil {
			return "", "", err
		}
	}

	toPrint[idx] = d
	return d.TypeNamespace, d.TypeName, nil
}

func printTypeDef(
	c *types.Context,
	def types.TypeDef,
	toPrint map[types.TypeDefOrRef]types.TypeDef,
) (string, error) {
	buf := strings.Builder{}

	fieldList, err := def.ResolveFieldList(c)
	if err != nil {
		return "", err
	}
	isNewType := len(fieldList) == 1

	buf.WriteString("type ")
	buf.WriteString(def.TypeName)
	buf.WriteByte(' ')
	if !isNewType {
		buf.WriteString("struct {\n")
	}
	for _, field := range fieldList {
		if field.Flags&0x0010 != 0 {
			continue // skip static fields
		}

		sig, err := field.Signature.Reader().Field(c)
		if err != nil {
			return "", err
		}
		_, fieldType, err := printName(c, sig.Field, toPrint)
		if err != nil {
			return "", err
		}
		if !isNewType {
			buf.WriteByte('\t')
			buf.WriteString(field.Name)
			buf.WriteByte(' ')
		}
		buf.WriteString(fieldType)
		buf.WriteByte('\n')
	}
	if !isNewType {
		buf.WriteString("}\n")
	}

	return buf.String(), nil
}

func printName(
	c *types.Context,
	e types.Element,
	toPrint map[types.TypeDefOrRef]types.TypeDef,
) (namespace, name string, err error) {
	switch e.Type.Kind {
	case types.ELEMENT_TYPE_U1:
		name = "uint8"
	case types.ELEMENT_TYPE_U2, types.ELEMENT_TYPE_CHAR:
		name = "uint16"
	case types.ELEMENT_TYPE_U4:
		name = "uint32"
	case types.ELEMENT_TYPE_U8:
		name = "uint64"
	case types.ELEMENT_TYPE_I1:
		name = "int8"
	case types.ELEMENT_TYPE_I2:
		name = "int16"
	case types.ELEMENT_TYPE_I4:
		name = "int32"
	case types.ELEMENT_TYPE_I8:
		name = "int64"
	case types.ELEMENT_TYPE_U:
		name = "uint"
	case types.ELEMENT_TYPE_I:
		name = "int"
	case types.ELEMENT_TYPE_R4:
		name = "float32"
	case types.ELEMENT_TYPE_R8:
		name = "float64"
	case types.ELEMENT_TYPE_BOOLEAN:
		name = "bool"
	case types.ELEMENT_TYPE_STRING:
		name = "string"
	case types.ELEMENT_TYPE_VALUETYPE, types.ELEMENT_TYPE_CLASS:
		namespace, name, err = queueTypeDefs(c, e, toPrint)
		if err != nil {
			return "", "", err
		}
	case types.ELEMENT_TYPE_VOID:
		if e.Pointers > 0 {
			name = "byte"
		}
	case types.ELEMENT_TYPE_OBJECT:
		name = "Object"
	case types.ELEMENT_TYPE_GENERICINST:
		_, genericName, err := c.ResolveTypeDefOrRefName(e.Type.TypeDef.Index)
		if err != nil {
			return "", "", err
		}
		if idx := strings.IndexByte(genericName, '`'); idx > 0 {
			genericName = genericName[:idx]
		}
		name = genericName
		name += "<"

		for i, arg := range e.Type.TypeDef.Generics {
			_, argName, err := printName(c, types.Element{
				Type: arg,
			}, toPrint)
			if err != nil {
				return "", "", err
			}
			name += argName

			if i != len(e.Type.TypeDef.Generics)-1 {
				name += ","
			}
		}
		name += ">"
	case types.ELEMENT_TYPE_ARRAY:
		ns, elemName, err := printName(c, *e.Type.Array.Elem, toPrint)
		if err != nil {
			return "", "", err
		}
		name = fmt.Sprintf("%s[%d]", elemName, e.Type.Array.Size)
		namespace = ns
	case types.ELEMENT_TYPE_SZARRAY:
		ns, elemName, err := printName(c, *e.Type.SZArray.Elem, toPrint)
		if err != nil {
			return "", "", err
		}
		name = elemName + "[]"
		namespace = ns
	default:
		return "", "", fmt.Errorf("unexpected type %v", e.Type)
	}

	if e.IsConst {
		name = "/* const */" + name
	}

	for i := 0; i < e.Pointers; i++ {
		name = "*" + name
	}
	return namespace, name, nil
}

func collectParamNames(c *types.Context, def types.MethodDef) (map[int]string, error) {
	paramNames := map[int]string{}

	params, err := def.ResolveParamList(c)
	if err != nil {
		return nil, err
	}

	for _, param := range params {
		paramNames[int(param.Sequence)-1] = param.Name
	}
	return paramNames, nil
}

func printMethod(
	c *types.Context,
	def types.MethodDef,
	toPrint map[types.TypeDefOrRef]types.TypeDef,
) (string, error) {
	r := def.Signature.Reader()

	method, err := r.Method(c)
	if err != nil {
		return "", err
	}
	count := len(method.Params)

	paramNames, err := collectParamNames(c, def)
	if err != nil {
		return "", err
	}

	log := strings.Builder{}
	log.WriteString("func ")
	log.WriteString(def.Name)
	log.WriteByte('(')
	if count > 0 {
		log.WriteByte('\n')
		for i := 0; i < count; i++ {
			log.WriteByte('\t')
			if paramName, ok := paramNames[i]; ok && paramName != "" {
				log.WriteString(paramName)
			} else {
				log.WriteString("p")
				log.WriteString(strconv.Itoa(i))
			}

			log.WriteByte(' ')

			_, typeName, err := printName(c, method.Params[i], toPrint)
			if err != nil {
				return "", err
			}
			log.WriteString(typeName)
			log.WriteString(",\n")
		}
	}
	log.WriteString(") ")

	_, typeName, err := printName(c, method.Return, toPrint)
	if err != nil {
		return "", err
	}
	log.WriteString(typeName)
	log.WriteByte('\n')

	return log.String(), nil
}
