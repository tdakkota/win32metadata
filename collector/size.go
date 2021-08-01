package collector

import (
	"fmt"

	"github.com/tdakkota/win32metadata/types"
)

func (c *collector) computeSize(e types.Element, ptrSize int) (int, error) {
	if e.Pointers > 0 {
		return ptrSize, nil
	}

	size := 0
	switch e.Type.Kind {
	case types.ELEMENT_TYPE_U1, types.ELEMENT_TYPE_I1:
		size = 1
	case types.ELEMENT_TYPE_U2, types.ELEMENT_TYPE_I2, types.ELEMENT_TYPE_CHAR:
		size = 2
	case types.ELEMENT_TYPE_U4, types.ELEMENT_TYPE_I4, types.ELEMENT_TYPE_R4, types.ELEMENT_TYPE_BOOLEAN:
		size = 4
	case types.ELEMENT_TYPE_U8, types.ELEMENT_TYPE_I8, types.ELEMENT_TYPE_R8:
		size = 8
	case types.ELEMENT_TYPE_U, types.ELEMENT_TYPE_I, types.ELEMENT_TYPE_OBJECT:
		size = ptrSize
	case types.ELEMENT_TYPE_VALUETYPE, types.ELEMENT_TYPE_CLASS:
		idx := e.Type.TypeDef.Index
		if v, ok := c.typeSizes[idx]; ok {
			return v, nil
		}

		elemType, err := c.findNamedType(idx)
		if err != nil {
			return 0, fmt.Errorf("find type %d: %w", idx, err)
		}
		d := elemType.Def

		fields, err := d.ResolveFieldList(c.ctx)
		if err != nil {
			return 0, err
		}

		for _, field := range fields {
			if field.Flags.Static() {
				continue
			}

			sig, err := field.Signature.Reader().Field(c.ctx)
			if err != nil {
				return 0, err
			}

			fieldSize, err := c.computeSize(sig.Field, ptrSize)
			if err != nil {
				return 0, fmt.Errorf(
					"get size of %s.%s, field %s: %w",
					d.TypeNamespace, d.TypeName, field.Name, err,
				)
			}

			if elemType.Union && size < fieldSize {
				size = fieldSize
			} else {
				size += fieldSize
			}
		}
		c.typeSizes[idx] = size
	case types.ELEMENT_TYPE_VOID:
		size = 0
	default:
		return 0, fmt.Errorf("can't count size of %v", e)
	}

	return size, nil
}
