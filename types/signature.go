package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Reader creates SignatureReader.
func (s Signature) Reader() *SignatureReader {
	return &SignatureReader{
		sig: s,
	}
}

// SignatureReader is a helper to read Signature.
type SignatureReader struct {
	sig    Signature
	offset int
}

// Peek peeks unsigned integer from Signature blob, returns it and its size in bytes.
// If there is no data anymore, ok is false.
func (s *SignatureReader) Peek() (value uint32, size int, ok bool) {
	if s.offset >= len(s.sig) {
		return 0, 0, false
	}

	b := []byte(s.sig[s.offset:])
	if len(b) == 0 {
		return 0, 0, false
	}
	switch {
	case b[0]&0x80 == 0:
		return uint32(b[0]), 1, true
	case b[0]&0xC0 == 0x80:
		if len(b) < 2 {
			return 0, 0, false
		}
		return uint32(binary.BigEndian.Uint16([]byte{b[0] & 0x3F, b[1]})), 2, true
	default:
		if len(b) < 4 {
			return 0, 0, false
		}
		return binary.BigEndian.Uint32([]byte{b[0] & 0x1F, b[1], b[2], b[3]}), 4, true
	}
}

// NextIs peeks unsigned integer from Signature blob and compares it with given value.
// If values are equal, increases offset.
func (s *SignatureReader) NextIs(expect uint32) bool {
	value, size, ok := s.Peek()
	if !ok || value != expect {
		return false
	}

	s.offset += size
	return true
}

// Read reads unsigned integer from Signature blob and returns it.
// If there is no data anymore, ok is false.
func (s *SignatureReader) Read() (value uint32, ok bool) {
	var size int
	value, size, ok = s.Peek()
	if !ok {
		return
	}
	s.offset += size
	return
}

func (s *SignatureReader) modifiers() (result []TypeDefOrRef) {
	for {
		value, size, ok := s.Peek()
		if !ok || (value != 32 && value != 31) {
			break
		}
		s.offset += size
		result = append(result, TypeDefOrRef(value))
	}

	return
}

// Element represents one parameter or result in Signature.
type Element struct {
	Type     ElementType
	Pointers int
	ByRef    bool
	IsConst  bool
	IsArray  bool
}

func (s *SignatureReader) elementType(c *Context) (ElementType, error) {
	var t ElementType

	value, ok := s.Read()
	if !ok {
		return t, io.ErrUnexpectedEOF
	}

	if t.FromCode(value) {
		return t, nil
	}

	// TODO(tdakkota): complete implementation
	switch ElementTypeKind(value) {
	case ELEMENT_TYPE_VALUETYPE, ELEMENT_TYPE_CLASS:
		r, ok := s.Read()
		if !ok {
			return t, io.ErrUnexpectedEOF
		}
		t.TypeDef.Index = TypeDefOrRef(r)

		return t, nil
	case ELEMENT_TYPE_VAR:
		r, ok := s.Read()
		if !ok {
			return t, io.ErrUnexpectedEOF
		}
		t.GenericTypeVar.Index = r

		return t, nil
	case ELEMENT_TYPE_MVAR:
		r, ok := s.Read()
		if !ok {
			return t, io.ErrUnexpectedEOF
		}
		t.GenericMethodVar.Index = r

		return t, nil
	case ELEMENT_TYPE_ARRAY:
		elem, err := s.NextElement(c)
		if err != nil {
			return ElementType{}, err
		}
		// TODO(tdakkota): complete decoding according to II.23.2.13 ArrayShape.
		s.Read() // rank
		s.Read() // bounds count

		size, ok := s.Read()
		if !ok {
			return ElementType{}, io.ErrUnexpectedEOF
		}
		t.Array = ElementTypeArray{
			Elem: &elem,
			Size: size,
		}
		return t, nil
	case ELEMENT_TYPE_GENERICINST:
		s.Read() // (CLASS | VALUETYPE)
		r, ok := s.Read()
		if !ok {
			return t, io.ErrUnexpectedEOF
		}
		t.TypeDef.Index = TypeDefOrRef(r)

		args, ok := s.Read() // GenArgCount
		if !ok {
			return t, io.ErrUnexpectedEOF
		}
		for i := uint32(0); i < args; i++ {
			arg, err := s.elementType(c)
			if err != nil {
				return t, err
			}
			t.TypeDef.Generics = append(t.TypeDef.Generics, arg)
		}

		return t, nil
	case ELEMENT_TYPE_SZARRAY:
		elem, err := s.NextElement(c)
		if err != nil {
			return ElementType{}, err
		}

		t.SZArray = ElementTypeSZArray{
			Elem: &elem,
		}
		return t, nil
	default:
		return t, fmt.Errorf("unexpected element type %#x", value)
	}
}

func (s *SignatureReader) isConst(c *Context) (bool, error) {
	for _, mod := range s.modifiers() {
		namespace, name, err := c.ResolveTypeDefOrRefName(mod)
		if err != nil {
			return false, err
		}

		if namespace == "System.Runtime.CompilerServices" && name == "IsConst" {
			return true, nil
		}
	}

	return false, nil
}

// NextElement returns next Element in signature.
func (s *SignatureReader) NextElement(c *Context) (e Element, _ error) {
	isConst, err := s.isConst(c)
	if err != nil {
		return e, err
	}
	e.IsConst = isConst

	e.ByRef = s.NextIs(uint32(ELEMENT_TYPE_BYREF))
	if s.NextIs(uint32(ELEMENT_TYPE_VOID)) {
		e.Type = ElementType{
			Kind: ELEMENT_TYPE_VOID,
		}
		return e, nil
	}
	e.IsArray = s.NextIs(uint32(ELEMENT_TYPE_ARRAY))
	for s.NextIs(uint32(ELEMENT_TYPE_PTR)) {
		e.Pointers++
	}

	elementType, err := s.elementType(c)
	if err != nil {
		return e, err
	}
	e.Type = elementType

	return
}

// MethodSignature is a II.23.2.1 MethodDefSig or II.23.2.2 MethodRefSig representation.
type MethodSignature struct {
	Flags           uint32
	GenericArgCount uint32
	Return          Element
	Params          []Element
}

// Method reads MethodSignature from Signature blob.
func (s *SignatureReader) Method(file *Context) (MethodSignature, error) {
	flags, ok := s.Read()
	if !ok {
		return MethodSignature{}, io.ErrUnexpectedEOF
	}

	const METHOD_DEF_SIG_FLAGS_GENERIC = 0x10
	var genericArgCount uint32
	if flags&METHOD_DEF_SIG_FLAGS_GENERIC == METHOD_DEF_SIG_FLAGS_GENERIC {
		genericArgCount, ok = s.Read()
		if !ok {
			return MethodSignature{}, io.ErrUnexpectedEOF
		}
	}

	count, ok := s.Read()
	if !ok {
		return MethodSignature{}, io.ErrUnexpectedEOF
	}

	returnType, err := s.NextElement(file)
	if err != nil {
		return MethodSignature{}, err
	}

	var params []Element
	if count > 0 {
		params = make([]Element, 0, count)
		for i := 0; i < int(count); i++ {
			t, err := s.NextElement(file)
			if err != nil {
				return MethodSignature{}, err
			}
			params = append(params, t)
		}
	}

	return MethodSignature{
		Flags:           flags,
		GenericArgCount: genericArgCount,
		Return:          returnType,
		Params:          params,
	}, nil
}

// FieldSignature is a II.23.2.4 FieldSig representation.
type FieldSignature struct {
	Field Element
}

// Field reads FieldSignature from Signature blob.
func (s *SignatureReader) Field(file *Context) (FieldSignature, error) {
	typ, ok := s.Read()
	if !ok {
		return FieldSignature{}, io.ErrUnexpectedEOF
	}
	if typ != 0x6 {
		return FieldSignature{}, fmt.Errorf("unexpected field type %d", typ)
	}

	e, err := s.NextElement(file)
	if err != nil {
		return FieldSignature{}, err
	}

	return FieldSignature{
		Field: e,
	}, nil
}
