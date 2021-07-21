package types

// ElementTypeKind is a II.23.1.16 Element types used in signatures representation kind.
type ElementTypeKind uint16

//go:generate go run golang.org/x/tools/cmd/stringer -type=ElementTypeKind

const (
	// ELEMENT_TYPE_END constant.
	ELEMENT_TYPE_END ElementTypeKind = 0x00
	// ELEMENT_TYPE_VOID constant.
	ELEMENT_TYPE_VOID ElementTypeKind = 0x01
	// ELEMENT_TYPE_BOOLEAN constant.
	ELEMENT_TYPE_BOOLEAN ElementTypeKind = 0x02
	// ELEMENT_TYPE_CHAR constant.
	ELEMENT_TYPE_CHAR ElementTypeKind = 0x03
	// ELEMENT_TYPE_I1 constant.
	ELEMENT_TYPE_I1 ElementTypeKind = 0x04
	// ELEMENT_TYPE_U1 constant.
	ELEMENT_TYPE_U1 ElementTypeKind = 0x05
	// ELEMENT_TYPE_I2 constant.
	ELEMENT_TYPE_I2 ElementTypeKind = 0x06
	// ELEMENT_TYPE_U2 constant.
	ELEMENT_TYPE_U2 ElementTypeKind = 0x07
	// ELEMENT_TYPE_I4 constant.
	ELEMENT_TYPE_I4 ElementTypeKind = 0x08
	// ELEMENT_TYPE_U4 constant.
	ELEMENT_TYPE_U4 ElementTypeKind = 0x09
	// ELEMENT_TYPE_I8 constant.
	ELEMENT_TYPE_I8 ElementTypeKind = 0x0a
	// ELEMENT_TYPE_U8 constant.
	ELEMENT_TYPE_U8 ElementTypeKind = 0x0b
	// ELEMENT_TYPE_R4 constant.
	ELEMENT_TYPE_R4 ElementTypeKind = 0x0c
	// ELEMENT_TYPE_R8 constant.
	ELEMENT_TYPE_R8 ElementTypeKind = 0x0d
	// ELEMENT_TYPE_STRING constant.
	ELEMENT_TYPE_STRING ElementTypeKind = 0x0e
	// ELEMENT_TYPE_PTR constant.
	ELEMENT_TYPE_PTR ElementTypeKind = 0x0f
	// ELEMENT_TYPE_BYREF constant.
	ELEMENT_TYPE_BYREF ElementTypeKind = 0x10
	// ELEMENT_TYPE_VALUETYPE constant.
	ELEMENT_TYPE_VALUETYPE ElementTypeKind = 0x11
	// ELEMENT_TYPE_CLASS constant.
	ELEMENT_TYPE_CLASS ElementTypeKind = 0x12
	// ELEMENT_TYPE_VAR constant.
	ELEMENT_TYPE_VAR ElementTypeKind = 0x13
	// ELEMENT_TYPE_ARRAY constant.
	ELEMENT_TYPE_ARRAY ElementTypeKind = 0x14
	// ELEMENT_TYPE_GENERICINST constant.
	ELEMENT_TYPE_GENERICINST ElementTypeKind = 0x15
	// ELEMENT_TYPE_TYPEDBYREF constant.
	ELEMENT_TYPE_TYPEDBYREF ElementTypeKind = 0x16
	// ELEMENT_TYPE_I constant.
	ELEMENT_TYPE_I ElementTypeKind = 0x18
	// ELEMENT_TYPE_U constant.
	ELEMENT_TYPE_U ElementTypeKind = 0x19
	// ELEMENT_TYPE_FNPTR constant.
	ELEMENT_TYPE_FNPTR ElementTypeKind = 0x1b
	// ELEMENT_TYPE_OBJECT constant.
	ELEMENT_TYPE_OBJECT ElementTypeKind = 0x1c
	// ELEMENT_TYPE_SZARRAY constant.
	ELEMENT_TYPE_SZARRAY ElementTypeKind = 0x1d
	// ELEMENT_TYPE_MVAR constant.
	ELEMENT_TYPE_MVAR ElementTypeKind = 0x1e
	// ELEMENT_TYPE_CMOD_REQD constant.
	ELEMENT_TYPE_CMOD_REQD ElementTypeKind = 0x1f
	// ELEMENT_TYPE_CMOD_OPT constant.
	ELEMENT_TYPE_CMOD_OPT ElementTypeKind = 0x20
	// ELEMENT_TYPE_INTERNAL constant.
	ELEMENT_TYPE_INTERNAL ElementTypeKind = 0x21
	// ELEMENT_TYPE_MODIFIER constant.
	ELEMENT_TYPE_MODIFIER ElementTypeKind = 0x40
	// ELEMENT_TYPE_SENTINEL constant.
	ELEMENT_TYPE_SENTINEL ElementTypeKind = 0x41
	// ELEMENT_TYPE_PINNED constant.
	ELEMENT_TYPE_PINNED ElementTypeKind = 0x45
)

type ElementTypeArray struct {
	Elem *Element
	Size uint32
}

// ElementType is a II.23.1.16 Element types used in signatures representation kind.
type ElementType struct {
	Kind         ElementTypeKind
	GenericParam Index `table:"GenericParam"`
	Array        ElementTypeArray
	MethodDef    Index        `table:"MethodDef"`
	Field        Index        `table:"Field"`
	TypeDef      TypeDefOrRef `table:"TypeDef"`
}

// FromCode tries to map code to ElementTypeKind.
func (e *ElementType) FromCode(code uint32) bool {
	switch code {
	case 0x01:
		e.Kind = ELEMENT_TYPE_VOID
	case 0x02:
		e.Kind = ELEMENT_TYPE_BOOLEAN
	case 0x03:
		e.Kind = ELEMENT_TYPE_CHAR
	case 0x04:
		e.Kind = ELEMENT_TYPE_I1
	case 0x05:
		e.Kind = ELEMENT_TYPE_U1
	case 0x06:
		e.Kind = ELEMENT_TYPE_I2
	case 0x07:
		e.Kind = ELEMENT_TYPE_U2
	case 0x08:
		e.Kind = ELEMENT_TYPE_I4
	case 0x09:
		e.Kind = ELEMENT_TYPE_U4
	case 0x0a:
		e.Kind = ELEMENT_TYPE_I8
	case 0x0b:
		e.Kind = ELEMENT_TYPE_U8
	case 0x0c:
		e.Kind = ELEMENT_TYPE_R4
	case 0x0d:
		e.Kind = ELEMENT_TYPE_R8
	case 0x18:
		e.Kind = ELEMENT_TYPE_I
	case 0x19:
		e.Kind = ELEMENT_TYPE_U
	case 0x0e:
		e.Kind = ELEMENT_TYPE_STRING
	case 0x1c:
		e.Kind = ELEMENT_TYPE_OBJECT
	default:
		e.Kind = ElementTypeKind(code)
		return false
	}

	return true
}
