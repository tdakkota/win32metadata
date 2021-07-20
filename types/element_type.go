package types

// ElementType is a II.23.1.16 Element types used in signatures representation.
type ElementType uint16

const (
	// ELEMENT_TYPE_END constant.
	ELEMENT_TYPE_END ElementType = 0x00
	// ELEMENT_TYPE_VOID constant.
	ELEMENT_TYPE_VOID ElementType = 0x01
	// ELEMENT_TYPE_BOOLEAN constant.
	ELEMENT_TYPE_BOOLEAN ElementType = 0x02
	// ELEMENT_TYPE_CHAR constant.
	ELEMENT_TYPE_CHAR ElementType = 0x03
	// ELEMENT_TYPE_I1 constant.
	ELEMENT_TYPE_I1 ElementType = 0x04
	// ELEMENT_TYPE_U1 constant.
	ELEMENT_TYPE_U1 ElementType = 0x05
	// ELEMENT_TYPE_I2 constant.
	ELEMENT_TYPE_I2 ElementType = 0x06
	// ELEMENT_TYPE_U2 constant.
	ELEMENT_TYPE_U2 ElementType = 0x07
	// ELEMENT_TYPE_I4 constant.
	ELEMENT_TYPE_I4 ElementType = 0x08
	// ELEMENT_TYPE_U4 constant.
	ELEMENT_TYPE_U4 ElementType = 0x09
	// ELEMENT_TYPE_I8 constant.
	ELEMENT_TYPE_I8 ElementType = 0x0a
	// ELEMENT_TYPE_U8 constant.
	ELEMENT_TYPE_U8 ElementType = 0x0b
	// ELEMENT_TYPE_R4 constant.
	ELEMENT_TYPE_R4 ElementType = 0x0c
	// ELEMENT_TYPE_R8 constant.
	ELEMENT_TYPE_R8 ElementType = 0x0d
	// ELEMENT_TYPE_STRING constant.
	ELEMENT_TYPE_STRING ElementType = 0x0e
	// ELEMENT_TYPE_PTR constant.
	ELEMENT_TYPE_PTR ElementType = 0x0f
	// ELEMENT_TYPE_BYREF constant.
	ELEMENT_TYPE_BYREF ElementType = 0x10
	// ELEMENT_TYPE_VALUETYPE constant.
	ELEMENT_TYPE_VALUETYPE ElementType = 0x11
	// ELEMENT_TYPE_CLASS constant.
	ELEMENT_TYPE_CLASS ElementType = 0x12
	// ELEMENT_TYPE_VAR constant.
	ELEMENT_TYPE_VAR ElementType = 0x13
	// ELEMENT_TYPE_ARRAY constant.
	ELEMENT_TYPE_ARRAY ElementType = 0x14
	// ELEMENT_TYPE_GENERICINST constant.
	ELEMENT_TYPE_GENERICINST ElementType = 0x15
	// ELEMENT_TYPE_TYPEDBYREF constant.
	ELEMENT_TYPE_TYPEDBYREF ElementType = 0x16
	// ELEMENT_TYPE_I constant.
	ELEMENT_TYPE_I ElementType = 0x18
	// ELEMENT_TYPE_U constant.
	ELEMENT_TYPE_U ElementType = 0x19
	// ELEMENT_TYPE_FNPTR constant.
	ELEMENT_TYPE_FNPTR ElementType = 0x1b
	// ELEMENT_TYPE_OBJECT constant.
	ELEMENT_TYPE_OBJECT ElementType = 0x1c
	// ELEMENT_TYPE_SZARRAY constant.
	ELEMENT_TYPE_SZARRAY ElementType = 0x1d
	// ELEMENT_TYPE_MVAR constant.
	ELEMENT_TYPE_MVAR ElementType = 0x1e
	// ELEMENT_TYPE_CMOD_REQD constant.
	ELEMENT_TYPE_CMOD_REQD ElementType = 0x1f
	// ELEMENT_TYPE_CMOD_OPT constant.
	ELEMENT_TYPE_CMOD_OPT ElementType = 0x20
	// ELEMENT_TYPE_INTERNAL constant.
	ELEMENT_TYPE_INTERNAL ElementType = 0x21
	// ELEMENT_TYPE_MODIFIER constant.
	ELEMENT_TYPE_MODIFIER ElementType = 0x40
	// ELEMENT_TYPE_SENTINEL constant.
	ELEMENT_TYPE_SENTINEL ElementType = 0x41
	// ELEMENT_TYPE_PINNED constant.
	ELEMENT_TYPE_PINNED ElementType = 0x45
)
