package collector

import (
	"fmt"

	"github.com/tdakkota/win32metadata/md"
	"github.com/tdakkota/win32metadata/types"
)

type namedType struct {
	Name      string
	Namespace string
	Def       types.TypeDef
	Union     bool
}

func (c *collector) findNamedType(idx types.TypeDefOrRef) (namedType, error) {
	if tableType, ok := idx.Table(); ok && tableType == md.TypeRef {
		row, ok := idx.Row(c.ctx)
		if !ok {
			return namedType{}, fmt.Errorf("TypeRef %d not found", idx.TableIndex())
		}

		var ref types.TypeRef
		if err := ref.FromRow(row); err != nil {
			return namedType{}, err
		}
		if ref.TypeNamespace == "System" && ref.TypeName == "Guid" {
			return namedType{
				Name:      ref.TypeName,
				Namespace: ref.TypeNamespace,
			}, nil
		}
	}

	d, err := c.resolveTypeDef(idx)
	if err != nil {
		return namedType{}, err
	}

	return namedType{
		Name:      d.TypeName,
		Namespace: d.TypeNamespace,
		Def:       d,
		Union:     d.Flags.ExplicitLayout(),
	}, nil
}

func (c *collector) goType(imp *imports, e types.Element) (string, error) {
	var name string
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
		elemType, err := c.findNamedType(e.Type.TypeDef.Index)
		if err != nil {
			return "", err
		}

		ns := elemType.Namespace
		n := elemType.Name
		switch {
		case elemType.Union:
			size, err := c.computeSize(e, 8)
			if err != nil {
				return "", err
			}
			name = fmt.Sprintf("[%d]byte", size)
		case ns == "System" && n == "Guid":
			name = "[128]byte"
		case ns == "Windows.Win32.System.WinRT" && n == "HSTRING":
			name = "[]uint16"
		case ns == "Windows.Win32.System.WinRT" && n == "IInspectable":
			name = "uintptr"
		default:
			name = imp.Import(elemType)
		}
	case types.ELEMENT_TYPE_VOID:
		if e.Pointers > 0 {
			name = "byte"
		}
	case types.ELEMENT_TYPE_OBJECT:
		name = "uintptr"
	case types.ELEMENT_TYPE_GENERICINST:
		return "", fmt.Errorf("generic types are not supported yet")
	case types.ELEMENT_TYPE_ARRAY:
		elemType, err := c.goType(imp, *e.Type.Array.Elem)
		if err != nil {
			return "", err
		}
		name = fmt.Sprintf("[%d]%s", e.Type.Array.Size, elemType)
	case types.ELEMENT_TYPE_SZARRAY:
		elemType, err := c.goType(imp, *e.Type.SZArray.Elem)
		if err != nil {
			return "", err
		}
		name = "[]" + elemType
	default:
		return "", fmt.Errorf("unexpected type %v", e.Type)
	}

	if e.IsConst {
		name = "/* const */" + name
	}

	for i := 0; i < e.Pointers; i++ {
		name = "*" + name
	}
	return name, nil
}
