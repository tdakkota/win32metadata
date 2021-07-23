package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignatureReader_Method(t *testing.T) {
	// From Windows.Win32.winmd, may be various depending on file
	// Does not matter for testing purposes.
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
				Convention: DefaultCallingConvention,
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
				Convention: DefaultCallingConvention,
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
				Convention: DefaultCallingConvention,
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
				Convention: HasThisCallingConvention,
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
				Convention: HasThisCallingConvention,
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
