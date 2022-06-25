package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignatureReader_Method(t *testing.T) {
	// TypeDefOrRef may be various depending on file, but it
	// does not matter for testing purposes.
	var (
		NTSTATUS = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(989)},
		}
		HANDLE = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(181)},
		}
		BOOL = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(89)},
		}
		DWORD = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(997)},
		}
		KEY_VALUE_ENTRY = ElementType{
			Kind:    ELEMENT_TYPE_VALUETYPE,
			TypeDef: ElementTypeTypeDef{Index: TypeDefOrRef(4905)},
		}
		IPropertySet = ElementType{
			Kind:    ELEMENT_TYPE_CLASS,
			TypeDef: ElementTypeTypeDef{Index: 15285},
		}
		ResourceCandidate = ElementType{
			Kind: ELEMENT_TYPE_CLASS,
			TypeDef: ElementTypeTypeDef{
				Index: 6077,
			},
		}
		mvar = func(idx uint32) ElementType {
			return ElementType{
				Kind: ELEMENT_TYPE_MVAR,
				GenericMethodVar: ElementTypeGenericMethodVar{
					Index: idx,
				},
			}
		}
	)

	tests := []struct {
		name   string
		sig    Signature
		expect MethodSignature
	}{
		{
			"RtlNtStatusToDosError",
			Signature{
				0,            // Calling convention
				1,            // Parameter count
				9,            // Return type
				17, 131, 221, // First param definition (ELEMENT_TYPE_VALUETYPE + TypeDefOrRef index)
			},
			MethodSignature{
				Flags: 0,
				Return: Element{
					Type: ElementType{Kind: ELEMENT_TYPE_U4},
				},
				Params: []Element{
					{Type: NTSTATUS},
				},
			},
		},
		{
			"DuplicateHandle",
			Signature{
				0,
				7,
				17, 89, // BOOL
				17, 128, 181, // HANDLE
				17, 128, 181, // HANDLE
				17, 128, 181, // HANDLE
				15, 17, 128, 181, // *HANDLE = LPHANDLE
				9,      // uint32
				17, 89, // BOOL
				17, 131, 229, // DWORD
			},
			MethodSignature{
				Flags: 0,
				Return: Element{
					Type: BOOL,
				},
				Params: []Element{
					{Type: HANDLE},
					{Type: HANDLE},
					{Type: HANDLE},
					{Type: HANDLE, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}},
					{Type: BOOL},
					{Type: DWORD},
				},
			},
		},
		{
			"NtQueryMultipleValueKey",
			Signature{
				0,
				6,
				17, 131, 221, // NTSTATUS
				17, 128, 181, // HANDLE
				15, 17, 147, 41, // *KEY_VALUE_ENTRY
				9,     // uint32
				15, 1, // *void
				15, 9, // *uint32
				15, 9, // *uint32
			},
			MethodSignature{
				Flags: 0,
				Return: Element{
					Type: NTSTATUS,
				},
				Params: []Element{
					{Type: HANDLE},
					{Type: KEY_VALUE_ENTRY, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}},
					{Type: ElementType{Kind: ELEMENT_TYPE_VOID}, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}, Pointers: 1},
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}, Pointers: 1},
				},
			},
		},
		{
			"AppExtension.GetExtensionPropertiesAsync",
			Signature{
				32,               // HASTHIS calling convention, see II.23.2.1 MethodDefSig
				0,                // No parameters
				21, 18, 188, 181, // Generic result type IAsyncOperation
				1,            // One generic argument
				18, 187, 181, // Argument is a class IPropertySet
			},
			MethodSignature{
				Flags: 32,
				Return: Element{Type: ElementType{
					Kind: ELEMENT_TYPE_GENERICINST,
					TypeDef: ElementTypeTypeDef{
						Index:    15541,
						Generics: []ElementType{IPropertySet},
					},
				}},
			},
		},
		{
			"ResourceCandidateVectorView.GetMany",
			Signature{
				32,
				2,
				9,
				9,
				29,           // SZARRAY argument of
				18, 151, 189, // ELEMENT_TYPE_CLASS ResourceCandidate
			},
			MethodSignature{
				Flags: 32,
				Return: Element{
					Type: ElementType{Kind: ELEMENT_TYPE_U4},
				},
				Params: []Element{
					{Type: ElementType{Kind: ELEMENT_TYPE_U4}},
					{
						Type: ElementType{
							Kind: ELEMENT_TYPE_SZARRAY,
							SZArray: ElementTypeSZArray{Elem: &Element{
								Type: ResourceCandidate,
							}},
						},
					},
				},
			},
		},
		{
			// VI.B.4.3 Metadata
			// 	class Phone<K,V> {
			// 	...
			// 	static void AddOne<KK,VV>(Phone<KK,VV> phone, KK kk, VV vv) { // reading this signature
			"GenericMethodExample",
			Signature{
				0x10, // IMAGE_CEE_CS_CALLCONV_GENERIC
				0x02, // GenParamCount = 2 (2 generic parameters for this method: KK and VV
				0x03, // ParamCount = 3 (phone, kk and vv)
				0x01, // RetType = ELEMENT_TYPE_VOID
				0x15, // Param-0: ELEMENT_TYPE_GENERICINST
				0x12, // 	ELEMENT_TYPE_CLASS
				0x08, // 	TypeDefOrRef coded index for class "Phone<KK,VV>"
				0x02, // 	GenArgCount = 2
				0x1e, // 	ELEMENT_TYPE_MVAR
				0x00, // 	!!0 (KK in AddOne<KK,VV>)
				0x1e, // 	ELEMENT_TYPE_MVAR
				0x01, // 	!!1 (VV in AddOne<KK,VV>)
				0x1e, // Param-1 ELEMENT_TYPE_MVAR
				0x00, // !!0 (KK in AddOne<KK,VV>)
				0x1e, // Param-2 ELEMENT_TYPE_MVAR
				0x01, // !!1 (VV in AddOne<KK,VV>)
			},
			MethodSignature{
				Flags:           0x10,
				GenericArgCount: 2,
				Return: Element{
					Type: ElementType{Kind: ELEMENT_TYPE_VOID},
				},
				Params: []Element{
					// Phone<KK,VV> phone,
					{
						Type: ElementType{
							Kind: ELEMENT_TYPE_GENERICINST,
							TypeDef: ElementTypeTypeDef{
								Index: 8,
								Generics: []ElementType{
									mvar(0),
									mvar(1),
								},
							},
						},
					},
					{Type: mvar(0)},
					{Type: mvar(1)},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := require.New(t)
			r := test.sig.Reader()

			method, err := r.Method(nil) // Context must not be needed anyway
			a.NoError(err)

			a.Equal(test.expect, method)
		})
	}
}

func TestSignatureReader_Field(t *testing.T) {
	tests := []struct {
		name   string
		sig    Signature
		expect FieldSignature
	}{
		{
			// CONSOLE_MODE is a enum, value__ is a enum value.
			"CONSOLE_MODE.value__",
			Signature{
				6, // Denotes that signature is a field type.
				9, // uint32
			},
			FieldSignature{Field: Element{
				Type: ElementType{Kind: ELEMENT_TYPE_U4},
			}},
		},
		{
			// CONSOLE_MODE is a enum, ENABLE_LINE_INPUT is a associated value.
			"CONSOLE_MODE.ENABLE_LINE_INPUT",
			Signature{
				6,  // Denotes that signature is a field type.
				17, // ELEMENT_TYPE_VALUETYPE
				5,  // TypeDefOrRef index
			},
			FieldSignature{Field: Element{
				Type: ElementType{
					Kind: ELEMENT_TYPE_VALUETYPE,
					TypeDef: ElementTypeTypeDef{
						Index: 5,
					},
				},
			}},
		},
		{
			// VI.B.4.3 Metadata
			// 	class Phone<K,V> {
			// 		private V[] vals; // reading this signature
			"ELEMENT_TYPE_VAR",
			Signature{
				0x06, // FIELD
				0x1D, // ELEMENT_TYPE_SZARRAY
				0x13, // ELEMENT_TYPE_VAR
				0x01, // 1, representing generic argument number 1 (i.e., "V")
			},
			FieldSignature{Field: Element{
				Type: ElementType{
					Kind: ELEMENT_TYPE_SZARRAY,
					SZArray: ElementTypeSZArray{
						Elem: &Element{
							Type: ElementType{
								Kind: ELEMENT_TYPE_VAR,
								GenericTypeVar: ElementTypeGenericTypeVar{
									Index: 1,
								},
							},
						},
					},
				},
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := require.New(t)
			r := test.sig.Reader()

			method, err := r.Field(nil) // Context must not be needed anyway
			a.NoError(err)

			a.Equal(test.expect, method)
		})
	}
}
